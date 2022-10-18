package main

import (
	"bufio"
	"fmt"
	"os"
	TDACola "rerepolez/cola"
	"rerepolez/errores"
	"rerepolez/votos"
	"strconv"
	"strings"
)

var cargos = [3]string{"Presidente", "Gobernador", "Intendente"}

var (
	Presidente string = cargos[0]
	Gobernador string = cargos[1]
	Intendente string = cargos[2]
)

func votacion(listaVotantes []votos.Votante, listaPartidos []votos.Partido, listaBlanco votos.Partido, votosImpugnados *int) {
	colaVotantes := TDACola.CrearColaEnlazada[votos.Votante]()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {

		input := strings.Split(scanner.Text(), " ")
		accion := input[0]

		switch accion {

		case "ingresar":
			pos, err := ingresarVotante(input, listaVotantes)

			if err != nil {
				fmt.Fprintln(os.Stdout, err.Error())
			} else {
				colaVotantes.Encolar(listaVotantes[pos])
				fmt.Println("OK")
			}

		case "votar":
			err := votar(input, colaVotantes, listaPartidos)

			if err != nil {
				fmt.Fprintln(os.Stdout, err.Error())
			} else {
				fmt.Println("OK")
			}

		case "deshacer":
			err := deshacerVoto(colaVotantes)

			if err != nil {
				fmt.Fprintln(os.Stdout, err.Error())
			} else {
				fmt.Println("OK")
			}

		case "fin-votar":
			voto, err := finalizarVoto(colaVotantes)

			if err != nil {
				fmt.Fprintln(os.Stdout, err.Error())
			} else {
				colaVotantes.Desencolar()
				fmt.Println("OK")
				repartirVotos(voto, votosImpugnados, listaBlanco, listaPartidos)
			}

		default:
			fmt.Fprintln(os.Stdout, "Input incorrecto")
		}
	}

	if !colaVotantes.EstaVacia() {
		fmt.Fprintln(os.Stdout, errores.ErrorCiudadanosSinVotar{}.Error())
	}
}

func ingresarVotante(input []string, listaVotantes []votos.Votante) (int, error) {
	dni, err := strconv.Atoi(input[1])

	if err != nil || dni < 0 {
		return -1, errores.DNIError{}
	}

	pos := buscarVotante(listaVotantes, 0, len(listaVotantes)-1, dni)
	if pos == -1 {
		return -1, errores.DNIFueraPadron{}
	}

	return pos, nil
}

func votar(input []string, colaVotantes TDACola.Cola[votos.Votante], listaPartidos []votos.Partido) error {
	if colaVotantes.EstaVacia() {
		return errores.FilaVacia{}
	} else if len(input) < 3 {
		return errores.ErrorTipoVoto{}
	}

	puesto := input[1]
	if puesto != Presidente && puesto != Gobernador && puesto != Intendente {
		return errores.ErrorTipoVoto{}
	}

	lista, err := strconv.Atoi(input[2])
	cantPartidos := len(listaPartidos)
	if err != nil || lista > cantPartidos {
		return errores.ErrorAlternativaInvalida{}
	}

	var errorVotanteFraudulento error
	for i, cargo := range cargos {
		if puesto == cargo {
			errorVotanteFraudulento = colaVotantes.VerPrimero().Votar(i, lista)
		}
	}

	if errorVotanteFraudulento != nil {
		colaVotantes.Desencolar()
		return errorVotanteFraudulento
	}

	return nil
}

func deshacerVoto(colaVotantes TDACola.Cola[votos.Votante]) error {
	if colaVotantes.EstaVacia() {
		return errores.FilaVacia{}
	}

	newError := colaVotantes.VerPrimero().Deshacer()
	if newError != nil {
		if newError.Error() != "ERROR: Sin voto a deshacer" {
			colaVotantes.Desencolar()
		}
		return newError
	}

	return nil
}

func finalizarVoto(colaVotantes TDACola.Cola[votos.Votante]) (votos.Voto, error) {
	if colaVotantes.EstaVacia() {
		return votos.Voto{}, errores.FilaVacia{}
	}

	voto, newError := colaVotantes.VerPrimero().FinVoto()

	if newError != nil {
		return votos.Voto{}, newError
	}

	return voto, nil
}

func repartirVotos(voto votos.Voto, votosImpugnados *int, listaBlanco votos.Partido, listaPartidos []votos.Partido) {
	if voto.Impugnado {
		*votosImpugnados++
		return
	}
	for i, voto := range voto.VotoPorTipo {
		if voto == votos.VOTO_EN_BLANCO {
			listaBlanco.VotadoPara(i)
		} else {
			listaPartidos[voto-1].VotadoPara(i)
		}
	}
}

func finalizar(listaBlanco votos.Partido, listaPartidos []votos.Partido, votosImpugnados *int) {
	for i, cargo := range cargos {
		if i == 0 {
			fmt.Fprintf(os.Stdout, "%s:\n", cargo)
		} else {
			fmt.Fprintf(os.Stdout, "\n%s:\n", cargo)
		}
		fmt.Fprintln(os.Stdout, listaBlanco.ObtenerResultado(i))
		for _, partido := range listaPartidos {
			fmt.Fprintln(os.Stdout, partido.ObtenerResultado(i))
		}
	}

	if *votosImpugnados == 1 {
		fmt.Fprintf(os.Stdout, "\nVotos Impugnados: %d voto\n", *votosImpugnados)
	} else {
		fmt.Fprintf(os.Stdout, "\nVotos Impugnados: %d votos\n", *votosImpugnados)
	}
}
