package votos

import (
	"rerepolez/errores"
	TDAPila "rerepolez/pila"
)

const (
	VOTANDO    = 0
	FINALIZADO = 1
)

type votanteImplementacion struct {
	dni    int
	votos  TDAPila.Pila[[2]int]
	estado int
}

func CrearVotante(dni int) Votante {
	v := new(votanteImplementacion)
	v.dni = dni
	v.votos = TDAPila.CrearPilaDinamica[[2]int]()
	v.estado = VOTANDO
	return v
}

func (votante votanteImplementacion) LeerDNI() int {
	return votante.dni
}

func (votante *votanteImplementacion) Votar(tipo int, alternativa int) error {
	if votante.estado == FINALIZADO {
		votanteFraudulento := new(errores.ErrorVotanteFraudulento)
		votanteFraudulento.Dni = votante.dni
		//return votanteFraudulento.Error()
	}
	voto := [2]int{tipo, alternativa}
	votante.votos.Apilar(voto)
	return nil
}

func (votante *votanteImplementacion) Deshacer() error {
	if votante.votos.EstaVacia() {
		//return errores.ErrorNoHayVotosAnteriores{}.Error()
	}
	votante.votos.Desapilar()
	return nil
}

func (votante *votanteImplementacion) FinVoto() (Voto, error) {
	votos := [3]int{}
	for !votante.votos.EstaVacia() {
		voto := votante.votos.Desapilar()
		if votos[voto[0]] == 0 {
			votos[voto[0]] = voto[1]
		}
	}
	votante.estado = FINALIZADO
	return Voto{votos, false}, nil
}
