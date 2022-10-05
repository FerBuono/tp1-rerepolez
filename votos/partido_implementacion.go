package votos

import "fmt"

type partidoImplementacion struct {
	nombre      string
	presidente  string
	votosPresid int
	gobernador  string
	votosGober  int
	intendente  string
	votosIntend int
}

type partidoEnBlanco struct {
	votosPresid int
	votosGober  int
	votosIntend int
}

func CrearPartido(nombre string, candidatos [CANT_VOTACION]string) Partido {
	partido := new(partidoImplementacion)
	partido.nombre = nombre
	partido.presidente = candidatos[PRESIDENTE]
	partido.votosPresid = 0
	partido.gobernador = candidatos[GOBERNADOR]
	partido.votosGober = 0
	partido.intendente = candidatos[INTENDENTE]
	partido.votosIntend = 0
	return partido
}

func CrearVotosEnBlanco() Partido {
	blanco := new(partidoEnBlanco)
	blanco.votosPresid = 0
	blanco.votosGober = 0
	blanco.votosIntend = 0
	return blanco
}

func (partido *partidoImplementacion) VotadoPara(tipo TipoVoto) {
	if tipo == "Presidente" {
		partido.votosPresid += 1
	} else if tipo == "Gobernador" {
		partido.votosGober += 1
	} else if tipo == "Intendente" {
		partido.votosIntend += 1
	}
}

func (partido partidoImplementacion) ObtenerResultado(tipo TipoVoto) string {
	if tipo == "Presidente" {
		return fmt.Sprintf("El Candidato %s tiene %i votos", partido.presidente, partido.votosPresid)
	} else if tipo == "Gobernador" {
		return fmt.Sprintf("El Candidato %s tiene %i votos", partido.gobernador, partido.votosGober)
	} else if tipo == "Intendente" {
		return fmt.Sprintf("El Candidato %s tiene %i votos", partido.intendente, partido.votosIntend)
	}
	return ""
}

func (blanco *partidoEnBlanco) VotadoPara(tipo TipoVoto) {
	if tipo == "Presidente" {
		blanco.votosPresid += 1
	} else if tipo == "Gobernador" {
		blanco.votosGober += 1
	} else if tipo == "Intendente" {
		blanco.votosIntend += 1
	}
}

func (blanco partidoEnBlanco) ObtenerResultado(tipo TipoVoto) string {
	if tipo == "Presidente" {
		return fmt.Sprintf("El Candidato %s tiene %i votos", partido.presidente, partido.votosPresid)
	} else if tipo == "Gobernador" {
		return fmt.Sprintf("El Candidato %s tiene %i votos", partido.gobernador, partido.votosGober)
	} else if tipo == "Intendente" {
		return fmt.Sprintf("El Candidato %s tiene %i votos", partido.intendente, partido.votosIntend)
	}
	return ""
}
