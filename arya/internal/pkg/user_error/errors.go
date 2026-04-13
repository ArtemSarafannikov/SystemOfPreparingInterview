package user_error

var (
	InternalError     Type = "Внутренняя ошибка"
	Unauthorized      Type = "Не авторизован"
	PermissionDenied  Type = "Нет доступа"
	NotFound          Type = "Не найдено"
	UserAlreadyExists Type = "Пользователь с таким никнеймом уже существует"
	WrongPassword     Type = "Неверный логин или пароль"
)
