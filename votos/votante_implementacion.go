package votos

import (
	"rerepolez/errores"
	TDAPila "rerepolez/pila"
)

const (
	VOTANDO = true
)

type voto struct {
	tipo        int
	alternativa int
}

type votanteImplementacion struct {
	dni     int
	votos   TDAPila.Pila[voto]
	votando bool
}

func CrearVotante(dni int) Votante {
	v := new(votanteImplementacion)
	v.dni = dni
	v.votos = TDAPila.CrearPilaDinamica[voto]()
	v.votando = VOTANDO
	return v
}

func (votante votanteImplementacion) LeerDNI() int {
	return votante.dni
}

func (votante *votanteImplementacion) Votar(tipo int, alternativa int) error {
	if !votante.votando {
		return votante.votanteFraudulento()
	}
	voto := voto{tipo, alternativa}
	votante.votos.Apilar(voto)
	return nil
}

func (votante *votanteImplementacion) Deshacer() error {
	if !votante.votando {
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
	if !votante.votando {
		return Voto{}, votante.votanteFraudulento()
	}

	votoFinal := [3]int{VOTO_EN_BLANCO, VOTO_EN_BLANCO, VOTO_EN_BLANCO}

	for !votante.votos.EstaVacia() {
		voto := votante.votos.Desapilar()

		if votoFinal[voto.tipo] == VOTO_EN_BLANCO || voto.alternativa == LISTA_IMPUGNA {
			votoFinal[voto.tipo] = voto.alternativa
		}
	}
	votante.votando = false
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
