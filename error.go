package irq

import "fmt"

type InterruptRequestIDAlreadyTakenError struct {
	ID uint
}

func (e InterruptRequestIDAlreadyTakenError) Error() string {
	return fmt.Sprintf("irq id is already taken: %d", e.ID)
}

type ReservedInterruptRequestIDError struct {
	ID uint
}

func (e ReservedInterruptRequestIDError) Error() string {
	return fmt.Sprintf("reserved irq id: %d", e.ID)
}

type InvalidContextError struct {
	Key string
}

func (e InvalidContextError) Error() string {
	return fmt.Sprintf("invalid context: %s key not found", e.Key)
}
