package main

import (
	"bufio"
	"fmt"
	"os"
	TDACola "rerepolez/cola"
	"rerepolez/errores"
	TDALista "rerepolez/lista-enlazada"
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

func Cantidad[T any](lista TDALista.Lista[T]) int {
	sum := 0
	lista.Iterar(func(dato T) bool {
		sum++
		return true
	})
	return sum
}

func guardarPartidos(partidos *os.File) TDALista.Lista[votos.Partido] {
	listaPartidos := TDALista.CrearListaEnlazada[votos.Partido]()

	scannerPartidos := bufio.NewScanner(partidos)

	nroLista := 1

	for scannerPartidos.Scan() {
		lista := strings.Split(scannerPartidos.Text(), ",")
		candidatos := [3]string{lista[1], lista[2], lista[3]}
		partido := votos.CrearPartido(nroLista, lista[0], candidatos)
		listaPartidos.InsertarUltimo(partido)
		nroLista++
	}
	return listaPartidos
}

func guardarPadron(padron *os.File) TDALista.Lista[votos.Votante] {
	listaVotantes := TDALista.CrearListaEnlazada[votos.Votante]()

	scannerPadron := bufio.NewScanner(padron)

	for scannerPadron.Scan() {
		dniVotante := scannerPadron.Text()
		dniVotanteNum, _ := strconv.Atoi(dniVotante)
		votante := votos.CrearVotante(dniVotanteNum)
		listaVotantes.InsertarUltimo(votante)
	}
	return listaVotantes
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
	listaBlanco := votos.CrearVotosEnBlanco()

	padron := AbrirArchivo(args[1])
	listaVotantes := guardarPadron(padron)
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
			}
			for iter := listaVotantes.Iterador(); iter.HaySiguiente(); {
				if dni != iter.VerActual().LeerDNI() {
					iter.Siguiente()
				} else {
					colaVotantes.Encolar(iter.VerActual())
					fmt.Println("OK")
					break
				}

				if !iter.HaySiguiente() {
					newError = new(errores.DNIFueraPadron)
					fmt.Fprintln(os.Stdout, newError.Error())
				}
			}

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
			cantPartidos := Cantidad(listaPartidos)
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
				break
			}

			fmt.Println(voto.VotoPorTipo)
			for i, voto := range voto.VotoPorTipo {
				if voto == votos.VOTO_EN_BLANCO {
					listaBlanco.VotadoPara(i)
					continue
				}
				for iter := listaPartidos.Iterador(); iter.HaySiguiente(); {
					if voto == iter.VerActual().LeerNroLista() {
						iter.VerActual().VotadoPara(i)
					}
					iter.Siguiente()
				}
			}

		default:
			fmt.Fprintln(os.Stdout, "Input incorrecto")
		}
	}
	fmt.Println("\nPresidente:")
	fmt.Fprintln(os.Stdout, listaBlanco.ObtenerResultado(votos.PRESIDENTE))
	for iter := listaPartidos.Iterador(); iter.HaySiguiente(); {
		fmt.Fprintln(os.Stdout, iter.VerActual().ObtenerResultado(votos.PRESIDENTE))
		iter.Siguiente()
	}
	fmt.Println("\nGobernador:")
	fmt.Fprintln(os.Stdout, listaBlanco.ObtenerResultado(votos.GOBERNADOR))
	for iter := listaPartidos.Iterador(); iter.HaySiguiente(); {
		fmt.Fprintln(os.Stdout, iter.VerActual().ObtenerResultado(votos.GOBERNADOR))
		iter.Siguiente()
	}
	fmt.Println("\nIntendente:")
	fmt.Fprintln(os.Stdout, listaBlanco.ObtenerResultado(votos.INTENDENTE))
	for iter := listaPartidos.Iterador(); iter.HaySiguiente(); {
		fmt.Fprintln(os.Stdout, iter.VerActual().ObtenerResultado(votos.INTENDENTE))
		iter.Siguiente()
	}
}
