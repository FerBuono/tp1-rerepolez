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
		return votante.votanteFraudulento()
	}

	voto := [2]int{int(tipo), alternativa}
	votante.votos.Apilar(voto)
	return nil
}

func (votante *votanteImplementacion) Deshacer() error {
	if votante.estado == FINALIZADO {
		return votante.votanteFraudulento()
	}
	if votante.votos.EstaVacia() {
		newError := new(errores.ErrorNoHayVotosAnteriores)
		return newError
	}
	votante.votos.Desapilar()
	return nil
}

func (votante *votanteImplementacion) FinVoto() (Voto, error) {
	if votante.estado == FINALIZADO {
		return Voto{}, votante.votanteFraudulento()
	}

	votoFinal := [3]int{VOTO_EN_BLANCO, VOTO_EN_BLANCO, VOTO_EN_BLANCO}

	for !votante.votos.EstaVacia() {
		voto := votante.votos.Desapilar()

		if votoFinal[voto[0]] == VOTO_EN_BLANCO || voto[1] == LISTA_IMPUGNA {
			votoFinal[voto[0]] = voto[1]
		}
	}
	votante.estado = FINALIZADO
	for _, voto := range votoFinal {
		if voto == LISTA_IMPUGNA {
			return Voto{votoFinal, true}, nil
		}
	}
	return Voto{votoFinal, false}, nil
}

func (votante *votanteImplementacion) votanteFraudulento() error {
	newError := new(errores.ErrorVotanteFraudulento)
	newError.Dni = votante.dni
	return newError
}
