package votos

type partidoImplementacion struct {
}

type partidoEnBlanco struct {
}

func CrearPartido(nombre string, candidatos [CANT_VOTACION]string) Partido {
	return nil
}

func CrearVotosEnBlanco() Partido {
	return nil
}

func (partido *partidoImplementacion) VotadoPara(tipo TipoVoto) {

}

func (partido partidoImplementacion) ObtenerResultado(tipo TipoVoto) string {
	return ""
}

func (blanco *partidoEnBlanco) VotadoPara(tipo TipoVoto) {

}

func (blanco partidoEnBlanco) ObtenerResultado(tipo TipoVoto) string {
	return ""
}
