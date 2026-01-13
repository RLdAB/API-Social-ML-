package domain

package domain

var (
    ErrUserNotFound     = errors.New("user not found")
    ErrSelfFollow       = errors.New("cannot follow yourself")
    ErrNotASeller       = errors.New("target user is not a seller")
    //Adicionar outros depois de orientações do Luiz
)