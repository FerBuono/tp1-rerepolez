package pila

const _CAPACIDAD_INICIAL = 10
const _VECES_A_AUMENTAR = 2
const _VECES_A_REDUCIR = 2
const _VALOR_PARA_REDUCIR = 4

/* Definición del struct pila proporcionado por la cátedra. */

type pilaDinamica[T any] struct {
	datos    []T
	cantidad int
}

func (p *pilaDinamica[T]) EstaVacia() bool {
	return p.cantidad == 0
}

func (p *pilaDinamica[T]) VerTope() T {
	if p.EstaVacia() {
		panic("La pila esta vacia")
	}
	return p.datos[p.cantidad-1]
}

func (p *pilaDinamica[T]) Apilar(elemento T) {
	if p.cantidad == cap(p.datos) {
		p.redimensionar(cap(p.datos) * _VECES_A_AUMENTAR)
	}

	p.datos[p.cantidad] = elemento
	p.cantidad++
}

func (p *pilaDinamica[T]) Desapilar() T {
	if p.EstaVacia() {
		panic("La pila esta vacia")
	}

	if p.cantidad <= cap(p.datos)/_VALOR_PARA_REDUCIR && cap(p.datos) > _CAPACIDAD_INICIAL {
		p.redimensionar(cap(p.datos) / _VECES_A_REDUCIR)
	}

	elemento := p.datos[p.cantidad-1]
	p.cantidad--

	return elemento
}

func (p *pilaDinamica[T]) redimensionar(nuevaCapacidad int) {
	nueva := make([]T, nuevaCapacidad)
	copy(nueva, p.datos)
	p.datos = nueva
}

func CrearPilaDinamica[T any]() Pila[T] {
	p := new(pilaDinamica[T])
	p.datos = make([]T, _CAPACIDAD_INICIAL)
	return p
}
