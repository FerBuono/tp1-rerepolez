package lista

type nodoLista[T any] struct {
	dato T
	prox *nodoLista[T]
}

type listaEnlazada[T any] struct {
	primero *nodoLista[T]
	ultimo  *nodoLista[T]
	largo   int
}

type iterListaEnlazada[T any] struct {
	lista    *listaEnlazada[T]
	actual   *nodoLista[T]
	anterior *nodoLista[T]
}

// Primitivas Lista

func (l *listaEnlazada[T]) EstaVacia() bool {
	return l.largo == 0 && l.primero == nil && l.ultimo == nil
}

func (l *listaEnlazada[T]) InsertarPrimero(dato T) {
	nuevo := l.crearNodo(dato)

	if l.EstaVacia() {
		l.ultimo = nuevo
	} else {
		nuevo.prox = l.primero
	}
	l.primero = nuevo

	l.largo++
}

func (l *listaEnlazada[T]) InsertarUltimo(dato T) {
	nuevo := l.crearNodo(dato)

	if l.EstaVacia() {
		l.primero = nuevo
	} else {
		l.ultimo.prox = nuevo
	}

	l.ultimo = nuevo
	l.largo++
}

func (l *listaEnlazada[T]) BorrarPrimero() T {
	if l.EstaVacia() {
		panic("La lista esta vacia")
	}

	primero := l.primero.dato

	if l.primero.prox == nil {
		l.ultimo = nil
	}
	l.primero = l.primero.prox

	l.largo--
	return primero
}

func (l *listaEnlazada[T]) VerPrimero() T {
	if l.EstaVacia() {
		panic("La lista esta vacia")
	}

	return l.primero.dato
}

func (l *listaEnlazada[T]) VerUltimo() T {
	if l.EstaVacia() {
		panic("La lista esta vacia")
	}

	return l.ultimo.dato
}

func (l *listaEnlazada[T]) Largo() int {
	return l.largo
}

func (l *listaEnlazada[T]) Iterador() IteradorLista[T] {
	iter := new(iterListaEnlazada[T])
	iter.lista = l
	iter.actual = l.primero
	iter.anterior = nil
	return iter
}

func (l *listaEnlazada[T]) Iterar(visitar func(T) bool) {
	actual := l.primero
	for actual != nil && visitar(actual.dato) {
		actual = actual.prox
	}
}

// Primitivas IteradorLista

func (i *iterListaEnlazada[T]) VerActual() T {
	if !i.HaySiguiente() {
		panic("El iterador termino de iterar")
	}
	return i.actual.dato
}

func (i *iterListaEnlazada[T]) HaySiguiente() bool {
	return i.actual != nil
}

func (i *iterListaEnlazada[T]) Siguiente() T {
	if !i.HaySiguiente() {
		panic("El iterador termino de iterar")
	}

	actual := i.actual.dato
	i.anterior = i.actual
	i.actual = i.actual.prox

	return actual
}

func (i *iterListaEnlazada[T]) Insertar(dato T) {
	nuevo := i.lista.crearNodo(dato)
	nuevo.prox = i.actual

	if i.anterior == nil {
		i.lista.primero = nuevo
		if i.actual == nil {
			i.lista.ultimo = nuevo
		}
	} else {
		i.anterior.prox = nuevo
		if i.actual == nil {
			i.lista.ultimo = nuevo
		}
	}
	i.actual = nuevo

	i.lista.largo++
}

func (i *iterListaEnlazada[T]) Borrar() T {
	if !i.HaySiguiente() {
		panic("El iterador termino de iterar")
	}

	dato := i.actual.dato

	if i.anterior == nil {
		i.lista.primero = i.actual.prox
		if i.actual.prox == nil {
			i.lista.ultimo = i.actual.prox
		}
	} else {
		i.anterior.prox = i.actual.prox
		if i.actual.prox == nil {
			i.lista.ultimo = i.anterior
		}
	}

	i.actual = i.actual.prox

	i.lista.largo--
	return dato
}

func (l *listaEnlazada[T]) crearNodo(dato T) *nodoLista[T] {
	nuevo := new(nodoLista[T])
	nuevo.dato = dato
	return nuevo
}

func CrearListaEnlazada[T any]() Lista[T] {
	l := new(listaEnlazada[T])
	return l
}
