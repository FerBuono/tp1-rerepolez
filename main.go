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

func guardarPadron(padron *os.File) []votos.Votante {
	listaVotantes := []votos.Votante{}

	scannerPadron := bufio.NewScanner(padron)
	for scannerPadron.Scan() {
		dniVotante, _ := strconv.Atoi(scannerPadron.Text())
		votante := votos.CrearVotante(dniVotante)
		listaVotantes = append(listaVotantes, votante)
	}
	return listaVotantes
}

func quickSort(lista []votos.Votante) []votos.Votante {
	if len(lista) < 2 {
		return lista
	}

	ini, fin := 0, len(lista)-1

	pivot := rand.Int() % len(lista)

	lista[pivot], lista[fin] = lista[fin], lista[pivot]

	for i := range lista {
		if lista[i].LeerDNI() < lista[fin].LeerDNI() {
			lista[ini], lista[i] = lista[i], lista[ini]
			ini++
		}
	}

	lista[ini], lista[fin] = lista[fin], lista[ini]

	quickSort(lista[:ini])
	quickSort(lista[ini+1:])

	return lista
}

func buscar(arr []votos.Votante, ini, fin, elemento int) int {
	if ini > fin {
		return -1
	}
	med := (ini + fin) / 2
	if arr[med].LeerDNI() == elemento {
		return med
	}
	if arr[med].LeerDNI() < elemento {
		return buscar(arr, med+1, fin, elemento)
	} else {
		return buscar(arr, ini, med-1, elemento)
	}

}

func votacion(listaVotantes []votos.Votante, listaPartidos []votos.Partido,
	listaBlanco votos.Partido, votosImpugnados int) {

	var newError error
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

			pos := buscar(listaVotantes, 0, len(listaVotantes)-1, dni)
			if pos == -1 {
				newError = new(errores.DNIFueraPadron)
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}

			colaVotantes.Encolar(listaVotantes[pos])
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
}

func finzalizar(listaBlanco votos.Partido, listaPartidos []votos.Partido, votosImpugnados int) {
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

	if votosImpugnados == 1 {
		fmt.Fprintf(os.Stdout, "\nVotos Impugnados: %d voto\n", votosImpugnados)
	} else {
		fmt.Fprintf(os.Stdout, "\nVotos Impugnados: %d votos\n", votosImpugnados)
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
	votacion(listaVotantes, listaPartidos, listaBlanco, votosImpugnados)
	finzalizar(listaBlanco, listaPartidos, votosImpugnados)
	/*
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

				pos := buscar(listaVotantes, 0, len(listaVotantes)-1, dni)
				if pos == -1 {
					newError = new(errores.DNIFueraPadron)
					fmt.Fprintln(os.Stdout, newError.Error())
					break
				}

				colaVotantes.Encolar(listaVotantes[pos])
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
		}*/
	/*
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

		if votosImpugnados == 1 {
			fmt.Fprintf(os.Stdout, "\nVotos Impugnados: %d voto\n", votosImpugnados)
		} else {
			fmt.Fprintf(os.Stdout, "\nVotos Impugnados: %d votos\n", votosImpugnados)
		}*/
}
