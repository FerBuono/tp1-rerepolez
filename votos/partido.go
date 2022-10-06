package votos

//Partido modela un partido político, con sus alternativas para cada uno de los tipos de votaciones
type Partido interface {

	//LeerNroLista devuelve el número de lista del partido
	LeerNroLista() int

	//VotadoPara indica que este Partido ha recibido un voto para el puesto indicado. Felicitaciones!
	VotadoPara(tipo int)

	//ObtenerResultado permite obtener el resultado de este Partido para el puesto indicado. El formato será el
	//conveniente para ser mostrado.
	ObtenerResultado(tipo int) string
}
