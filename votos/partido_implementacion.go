package votos

import "fmt"

type partidoImplementacion struct {
	nroLista    int
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

func CrearPartido(nombre string, candidatos [3]string) Partido {
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

func (partido *partidoImplementacion) VotadoPara(tipo int) {
	if tipo == PRESIDENTE {
		partido.votosPresid += 1
	} else if tipo == GOBERNADOR {
		partido.votosGober += 1
	} else if tipo == INTENDENTE {
		partido.votosIntend += 1
	}
}

func (partido partidoImplementacion) ObtenerResultado(tipo int) string {
	if tipo == PRESIDENTE {
		if partido.votosPresid == 1 {
			return fmt.Sprintf("%s - %s: %d voto", partido.nombre, partido.presidente, partido.votosPresid)
		}
		return fmt.Sprintf("%s - %s: %d votos", partido.nombre, partido.presidente, partido.votosPresid)
	} else if tipo == GOBERNADOR {
		if partido.votosGober == 1 {
			return fmt.Sprintf("%s - %s: %d voto", partido.nombre, partido.gobernador, partido.votosGober)
		}
		return fmt.Sprintf("%s - %s: %d votos", partido.nombre, partido.gobernador, partido.votosGober)
	} else if tipo == INTENDENTE {
		if partido.votosIntend == 1 {
			return fmt.Sprintf("%s - %s: %d voto", partido.nombre, partido.intendente, partido.votosIntend)
		}
		return fmt.Sprintf("%s - %s: %d votos", partido.nombre, partido.intendente, partido.votosIntend)
	}
	return ""
}

func (blanco *partidoEnBlanco) VotadoPara(tipo int) {
	if tipo == PRESIDENTE {
		blanco.votosPresid += 1
	} else if tipo == GOBERNADOR {
		blanco.votosGober += 1
	} else if tipo == INTENDENTE {
		blanco.votosIntend += 1
	}
}

func (blanco partidoEnBlanco) ObtenerResultado(tipo int) string {
	if tipo == PRESIDENTE {
		if blanco.votosPresid == 1 {
			return fmt.Sprintf("Votos en Blanco: %d voto", blanco.votosPresid)
		}
		return fmt.Sprintf("Votos en Blanco: %d votos", blanco.votosPresid)
	} else if tipo == GOBERNADOR {
		if blanco.votosGober == 1 {
			return fmt.Sprintf("Votos en Blanco: %d voto", blanco.votosGober)
		}
		return fmt.Sprintf("Votos en Blanco: %d votos", blanco.votosGober)
	} else if tipo == INTENDENTE {
		if blanco.votosIntend == 1 {
			return fmt.Sprintf("Votos en Blanco: %d voto", blanco.votosIntend)
		}
		return fmt.Sprintf("Votos en Blanco: %d votos", blanco.votosIntend)
	}
	return ""
}
