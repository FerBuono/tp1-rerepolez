package main

import (
	"bufio"
	"fmt"
	"os"
	"rerepolez/errores"
	"rerepolez/votos"
	"strconv"
	"strings"
)

const (
	_CANT_DIGITOS = 10
	_DIGITOS_DNI  = 8
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

func intPow(n, m int) int {
	if m == 0 {
		return 1
	}
	result := n
	for i := 2; i <= m; i++ {
		result *= n
	}
	return result
}

func countingPorCriterio(lista []votos.Votante, criterio func(int) int) []votos.Votante {
	frecuencias := make([]int, _CANT_DIGITOS)
	for _, votante := range lista {
		frecuencias[criterio(votante.LeerDNI())]++
	}

	frecSum := make([]int, _CANT_DIGITOS)
	sum := 0
	for i := range frecSum {
		frecSum[i] += sum
		sum += frecuencias[i]
	}

	arregloOrdenado := make([]votos.Votante, len(lista))
	for _, votante := range lista {
		arregloOrdenado[frecSum[criterio(votante.LeerDNI())]] = votante
		frecSum[criterio(votante.LeerDNI())]++
	}

	return arregloOrdenado
}

func ordenarListaVotantes(lista []votos.Votante) []votos.Votante {
	arregloOrdenado := lista
	for i := 0; i < _DIGITOS_DNI; i++ {
		arregloOrdenado = countingPorCriterio(arregloOrdenado, func(num int) int { return (num / intPow(10, i)) % 10 })
	}

	return arregloOrdenado
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
