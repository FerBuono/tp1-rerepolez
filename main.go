package main

import (
	"fmt"
	"os"
	"rerepolez/errores"
	"rerepolez/votos"
)

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
	var votosImpugnados *int = new(int)

	votacion(listaVotantes, listaPartidos, listaBlanco, votosImpugnados)
	finalizar(listaBlanco, listaPartidos, votosImpugnados)
}
