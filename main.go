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
	scannerPartidos := bufio.NewScanner(partidos)

	listaPartidos := TDALista.CrearListaEnlazada[votos.Partido]()

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
	var args = os.Args[1:]
	if len(args) < 2 {
		newError := new(errores.ErrorParametros)
		fmt.Fprintln(os.Stdout, newError.Error())
		os.Exit(0)
	}

	partidos := AbrirArchivo(args[0])
	listaPartidos := guardarPartidos(partidos)

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
				newError := new(errores.DNIError)
				fmt.Fprintln(os.Stdout, newError.Error())
			}
			for iter := listaVotantes.Iterador(); iter.HaySiguiente(); {
				if dni != iter.VerActual().LeerDNI() {
					iter.Siguiente()
				} else {
					colaVotantes.Encolar(iter.VerActual())
					break
				}

				if !iter.HaySiguiente() {
					newError := new(errores.DNIFueraPadron)
					fmt.Fprintln(os.Stdout, newError.Error())
				}
			}
			println("OK")

		case "votar":
			if colaVotantes.EstaVacia() {
				newError := new(errores.FilaVacia)
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			} else if len(input) < 3 {
				newError := new(errores.ErrorTipoVoto)
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}

			puesto := input[1]
			if puesto != "Presidente" && puesto != "Gobernador" && puesto != "Intendente" {
				newError := new(errores.ErrorTipoVoto)
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}

			lista, err := strconv.Atoi(input[2])
			cantPartidos := Cantidad(listaPartidos)
			if lista > cantPartidos || err != nil {
				newError := new(errores.ErrorAlternativaInvalida)
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}

			if puesto == "Presidente" {
				colaVotantes.VerPrimero().Votar(0, lista)
			}
			if puesto == "Gobernador" {
				colaVotantes.VerPrimero().Votar(1, lista)
			}
			if puesto == "Intendente" {
				colaVotantes.VerPrimero().Votar(2, lista)
			}
			println("OK")

		case "deshacer":
			if colaVotantes.EstaVacia() {
				newError := new(errores.FilaVacia)
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}
			colaVotantes.VerPrimero().Deshacer()
			println("OK")
		case "fin-votar":
			if colaVotantes.EstaVacia() {
				newError := new(errores.FilaVacia)
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}
		}
	}
}
