# Назначение
Отправка событий о звонках десктопного приложения платформы [1С-Connect](https://1c-connect.com/ru/) во внешние системы, например: CRM, ERP, PBX.

# Поддерживаемые системы:
* CRM
  * AmoCRM
* PBX
  * Asterisk (AMI)

Возможно список будет расширяться при наличии свободного времени и спроса на это дело.

# Cборка

```
git clone https://github.com/ros-tel/1c-connect-events.git 1c-connect-events
cd 1c-connect-events
go install
```

Для кроскомпиляции предварительно добавить в переменные окружения GOOS=windows и требуемую архитектуру (GOARCH)
```
export GOOS=windows
export GOARCH="amd64" (или "386")
```
