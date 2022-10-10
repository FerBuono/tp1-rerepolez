package votos

import "fmt"

type partidoImplementacion struct {
	nombre     string
	candidatos [3]string
	votos      [3]int
}

type partidoEnBlanco struct {
	votos [3]int
}

func CrearPartido(nombre string, candidatos [3]string) Partido {
	partido := new(partidoImplementacion)
	partido.nombre = nombre
	partido.candidatos = candidatos
	return partido
}

func CrearVotosEnBlanco() Partido {
	blanco := new(partidoEnBlanco)
	return blanco
}

func (partido *partidoImplementacion) VotadoPara(tipo int) {
	partido.votos[tipo]++
}

func (partido partidoImplementacion) ObtenerResultado(tipo int) string {
	if partido.votos[tipo] == 1 {
		return fmt.Sprintf("%s - %s: %d voto", partido.nombre, partido.candidatos[tipo], partido.votos[tipo])
	}
	return fmt.Sprintf("%s - %s: %d votos", partido.nombre, partido.candidatos[tipo], partido.votos[tipo])
}

func (blanco *partidoEnBlanco) VotadoPara(tipo int) {
	blanco.votos[tipo]++
}

func (blanco partidoEnBlanco) ObtenerResultado(tipo int) string {
	if blanco.votos[tipo] == 1 {
		return fmt.Sprintf("Votos en Blanco: %d voto", blanco.votos[tipo])
	}
	return fmt.Sprintf("Votos en Blanco: %d votos", blanco.votos[tipo])
}
