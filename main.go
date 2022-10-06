package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	TDACola "rerepolez/cola"
	"rerepolez/errores"
	"rerepolez/votos"
	"strconv"
	"strings"
)

func AbrirArchivo(archivo string) *os.File {
	file, err := os.Open(archivo)
	if err != nil {
		newError := new(errores.ErrorLeerArchivo)
		fmt.Fprintln(os.Stdout, newError.Error())
		os.Exit(0)
	}
	return file
}

func guardarPartidos(partidos *os.File) []votos.Partido {
	listaPartidos := []votos.Partido{}

	scannerPartidos := bufio.NewScanner(partidos)
	for scannerPartidos.Scan() {
		lista := strings.Split(scannerPartidos.Text(), ",")
		candidatos := [3]string{lista[1], lista[2], lista[3]}
		partido := votos.CrearPartido(lista[0], candidatos)
		listaPartidos = append(listaPartidos, partido)
	}
	return listaPartidos
}

func guardarPadron(padron *os.File) []int {
	listaVotantes := []int{}

	scannerPadron := bufio.NewScanner(padron)
	for scannerPadron.Scan() {
		dniVotante, _ := strconv.Atoi(scannerPadron.Text())
		listaVotantes = append(listaVotantes, dniVotante)
	}
	return listaVotantes
}

func quickSort(arr []int) []int {
	if len(arr) < 2 {
		return arr
	}

	ini, fin := 0, len(arr)-1

	pivot := rand.Int() % len(arr)

	arr[pivot], arr[fin] = arr[fin], arr[pivot]

	for i, _ := range arr {
		if arr[i] < arr[fin] {
			arr[ini], arr[i] = arr[i], arr[ini]
			ini++
		}
	}

	arr[ini], arr[fin] = arr[fin], arr[ini]

	quickSort(arr[:ini])
	quickSort(arr[ini+1:])

	return arr
}

func buscar(arr []int, ini, fin, elemento int) int {
	if ini > fin {
		return -1
	}
	med := (ini + fin) / 2
	if arr[med] == elemento {
		return med
	}
	if arr[med] < elemento {
		return buscar(arr, med+1, fin, elemento)
	} else {
		return buscar(arr, ini, med-1, elemento)
	}

}

func main() {
	var newError error

	var args = os.Args[1:]
	if len(args) < 2 {
		newError = new(errores.ErrorParametros)
		fmt.Fprintln(os.Stdout, newError.Error())
		os.Exit(0)
	}

	partidos := AbrirArchivo(args[0])
	listaPartidos := guardarPartidos(partidos)

	padron := AbrirArchivo(args[1])
	listaVotantes := quickSort(guardarPadron(padron))

	listaBlanco := votos.CrearVotosEnBlanco()
	votosImpugnados := 0

	colaVotantes := TDACola.CrearColaEnlazada[votos.Votante]()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {

		input := strings.Split(scanner.Text(), " ")
		accion := input[0]

		switch accion {

		case "ingresar":
			dni, err := strconv.Atoi(input[1])
			if err != nil || dni < 0 {
				newError = new(errores.DNIError)
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}

			if buscar(listaVotantes, 0, len(listaVotantes)-1, dni) == -1 {
				newError = new(errores.DNIFueraPadron)
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}

			votante := votos.CrearVotante(dni)
			colaVotantes.Encolar(votante)
			fmt.Println("OK")

		case "votar":
			if colaVotantes.EstaVacia() {
				newError = new(errores.FilaVacia)
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			} else if len(input) < 3 {
				newError = new(errores.ErrorTipoVoto)
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}

			puesto := input[1]
			if puesto != "Presidente" && puesto != "Gobernador" && puesto != "Intendente" {
				newError = new(errores.ErrorTipoVoto)
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}

			lista, err := strconv.Atoi(input[2])
			cantPartidos := len(listaPartidos)
			if lista > cantPartidos || err != nil {
				newError = new(errores.ErrorAlternativaInvalida)
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}

			if puesto == "Presidente" {
				newError = colaVotantes.VerPrimero().Votar(votos.PRESIDENTE, lista)
			}
			if puesto == "Gobernador" {
				newError = colaVotantes.VerPrimero().Votar(votos.GOBERNADOR, lista)
			}
			if puesto == "Intendente" {
				newError = colaVotantes.VerPrimero().Votar(votos.INTENDENTE, lista)
			}

			if newError != nil {
				fmt.Fprintln(os.Stdout, newError.Error())
				colaVotantes.Desencolar()
				break
			}

			fmt.Println("OK")

		case "deshacer":
			if colaVotantes.EstaVacia() {
				newError = new(errores.FilaVacia)
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}

			newError = colaVotantes.VerPrimero().Deshacer()
			if newError != nil {
				fmt.Fprintln(os.Stdout, newError.Error())
				if newError.Error() != "ERROR: Sin voto a deshacer" {
					colaVotantes.Desencolar()
				}
				break
			}

			fmt.Println("OK")

		case "fin-votar":
			if colaVotantes.EstaVacia() {
				newError = new(errores.FilaVacia)
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}

			voto, newError := colaVotantes.VerPrimero().FinVoto()

			if newError != nil {
				fmt.Fprintln(os.Stdout, newError.Error())
			} else {
				colaVotantes.Desencolar()
				fmt.Println("OK")
			}

			if voto.Impugnado {
				votosImpugnados++
				break
			}

			for i, voto := range voto.VotoPorTipo {
				if voto == votos.VOTO_EN_BLANCO {
					listaBlanco.VotadoPara(i)
					continue
				} else {
					listaPartidos[voto-1].VotadoPara(i)
				}
			}

		default:
			fmt.Fprintln(os.Stdout, "Input incorrecto")
		}
	}

	if !colaVotantes.EstaVacia() {
		newError = new(errores.ErrorCiudadanosSinVotar)
		fmt.Fprintln(os.Stdout, newError.Error())
	}

	fmt.Fprintln(os.Stdout, "Presidente:")
	fmt.Fprintln(os.Stdout, listaBlanco.ObtenerResultado(votos.PRESIDENTE))
	for _, partido := range listaPartidos {
		fmt.Fprintln(os.Stdout, partido.ObtenerResultado(votos.PRESIDENTE))
	}

	fmt.Fprintln(os.Stdout, "\nGobernador:")
	fmt.Fprintln(os.Stdout, listaBlanco.ObtenerResultado(votos.GOBERNADOR))
	for _, partido := range listaPartidos {
		fmt.Fprintln(os.Stdout, partido.ObtenerResultado(votos.GOBERNADOR))
	}

	fmt.Fprintln(os.Stdout, "\nIntendente:")
	fmt.Fprintln(os.Stdout, listaBlanco.ObtenerResultado(votos.INTENDENTE))
	for _, partido := range listaPartidos {
		fmt.Fprintln(os.Stdout, partido.ObtenerResultado(votos.INTENDENTE))
	}

	fmt.Fprintf(os.Stdout, "\nVotos Impugnados: %d votos\n", votosImpugnados)
}
