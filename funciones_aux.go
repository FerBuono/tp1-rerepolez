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

const (
	_CANT_DIGITOS = 10
)

func abrirArchivo(archivo string) *os.File {
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

func countingPorCriterio(lista []votos.Votante, criterio func(int) int) []votos.Votante {
	colas := make([]TDACola.Cola[votos.Votante], _CANT_DIGITOS)
	for i := range colas {
		colas[i] = TDACola.CrearColaEnlazada[votos.Votante]()
	}

	for _, votante := range lista {
		colas[criterio(votante.LeerDNI())].Encolar(votante)
	}

	listaOrdenada := make([]votos.Votante, len(lista))
	i := 0
	for _, cola := range colas {
		for !cola.EstaVacia() {
			listaOrdenada[i] = cola.Desencolar()
			i++
		}
	}

	return listaOrdenada
}

func ordenarListaVotantes(lista []votos.Votante) []votos.Votante {
	ordenadoPorPrimDig := countingPorCriterio(lista, func(num int) int { return num % 10 })
	ordenadoPorSegDig := countingPorCriterio(ordenadoPorPrimDig, func(num int) int { return (num / 10) % 10 })
	ordenadoPorTerDig := countingPorCriterio(ordenadoPorSegDig, func(num int) int { return (num / 100) % 10 })
	ordenadoPorCuarDig := countingPorCriterio(ordenadoPorTerDig, func(num int) int { return (num / 1000) % 10 })
	ordenadoPorQuinDig := countingPorCriterio(ordenadoPorCuarDig, func(num int) int { return (num / 10000) % 10 })
	ordenadoPorSextDig := countingPorCriterio(ordenadoPorQuinDig, func(num int) int { return (num / 100000) % 10 })
	ordenadoPorSeptDig := countingPorCriterio(ordenadoPorSextDig, func(num int) int { return (num / 1000000) % 10 })
	ordenadoPorOctDig := countingPorCriterio(ordenadoPorSeptDig, func(num int) int { return num / 10000000 })

	return ordenadoPorOctDig
}

func buscarVotante(arr []votos.Votante, ini, fin, elemento int) int {
	if ini > fin {
		return -1
	}
	med := (ini + fin) / 2
	if arr[med].LeerDNI() == elemento {
		return med
	}
	if arr[med].LeerDNI() < elemento {
		return buscarVotante(arr, med+1, fin, elemento)
	} else {
		return buscarVotante(arr, ini, med-1, elemento)
	}

}

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
	if puesto != "Presidente" && puesto != "Gobernador" && puesto != "Intendente" {
		return errores.ErrorTipoVoto{}
	}

	lista, err := strconv.Atoi(input[2])
	cantPartidos := len(listaPartidos)
	if err != nil || lista > cantPartidos {
		return errores.ErrorAlternativaInvalida{}
	}

	var errorVotanteFraudulento error
	if puesto == "Presidente" {
		errorVotanteFraudulento = colaVotantes.VerPrimero().Votar(votos.PRESIDENTE, lista)
	} else if puesto == "Gobernador" {
		errorVotanteFraudulento = colaVotantes.VerPrimero().Votar(votos.GOBERNADOR, lista)
	} else if puesto == "Intendente" {
		errorVotanteFraudulento = colaVotantes.VerPrimero().Votar(votos.INTENDENTE, lista)
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

	if *votosImpugnados == 1 {
		fmt.Fprintf(os.Stdout, "\nVotos Impugnados: %d voto\n", *votosImpugnados)
	} else {
		fmt.Fprintf(os.Stdout, "\nVotos Impugnados: %d votos\n", *votosImpugnados)
	}
}
