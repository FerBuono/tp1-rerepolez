package main

import (
	"fmt"
	"os"
	"rerepolez/errores"
	"rerepolez/votos"
)

func main() {

	var args = os.Args[1:]
	if len(args) < 2 {
		fmt.Fprintln(os.Stdout, errores.ErrorParametros{}.Error())
		os.Exit(0)
	}

	partidos := abrirArchivo(args[0])
	listaPartidos := guardarPartidos(partidos)
	partidos.Close()

	padron := abrirArchivo(args[1])
	listaVotantes := ordenarListaVotantes(guardarPadron(padron))
	padron.Close()

	listaBlanco := votos.CrearVotosEnBlanco()
	var votosImpugnados *int = new(int)

	votacion(listaVotantes, listaPartidos, listaBlanco, votosImpugnados)
	finalizar(listaBlanco, listaPartidos, votosImpugnados)
}
