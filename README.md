Biathlon Competition System Prototype
======================================

![biathlon](https://github.com/user-attachments/assets/2b57c0ca-26cf-4f29-98da-cf807dd10045)

$${\color{green}> \space \color{green}English \space \color{green}version}$$
---:
<br>
A prototype system for managing and analyzing biathlon competitions.
Processes configuration files and events, generates final reports.<br>

**Task Description:**
- [Описание задания на русском](README_TZ_rus.md)
- [Task description in English](README_TZ_eng.md)

---
## Installation
1. Ensure Go version 1.20 or newer is installed.
2. Clone the repository:
   - `git clone https://github.com/BiathlonRaceProto-Yadro.git`
   - `cd BiathlonRaceProto-Yadro`
3. Install dependencies:
   `make deps`

---
## Running via Makefile

**Commands:**
- `make build`        - Build the project
- `make run`          - Run with minimal output (as per assignment XD)
- `make run-fullOutput` - Run with full table report (default when invoking `make`)
- `make clean`        - Clean built files
- `make setup`        - Prepare directories and sample configs

**Examples:**
```
~\GolandProjects\BiathlonRaceProto-Yadro git:[main]
make run
"Сборка проекта..."
go build -o ./cmd/run/biathlon ./cmd/run/main.go
"Запуск приложения в minimal версии..."
./cmd/run/biathlon ./input/config/config.json ./input/events/events
Final Results:
. . .
```

---
## Manual Execution

1. Build:
   - Use `go build -o <executable_name> cmd/run/main.go`
   - Or skip and proceed to *step 2*<br><br>

   Example:
   ```
   ~\GolandProjects\BiathlonRaceProto-Yadro\cmd\run git:[main]
   go build -o biathlon main.go
   ```
2. Basic run:
   - `./biathlon <path_to_config.json> <path_to_events>`
   - Or, if skipped step 1: `go run main.go <path_to_config.json> <path_to_events>`<br><br>

   Example:
   ```
   ~\GolandProjects\BiathlonRaceProto-Yadro\cmd\run git:[main]
   go run main.go ..\..\input\config\config.json ..\..\input\events\events
   Final Results:
   . . .
   ```

3. Running with flags:
   - `-debug`       - Enable debug logs
   - `-info`        - INFO-level logs
   - `-error`       - ERROR-level logs
   - `-fullOutput`  - Full table report <br><br>

   Example:
   ```
   ~\GolandProjects\BiathlonRaceProto-Yadro\cmd\run git:[main]
   go run main.go -fullOutput -info ..\..\input\config\config.json ..\..\input\events\events
   {"time":"2025-05-07T03:19:51.2920861+03:00","level":"INFO","msg":"Участник зарегистрирован","time":"09:31:49.285","competitorID":3}
   . . 
   . . .
   {"time":"2025-05-07T03:19:51.3194668+03:00","level":"INFO","msg":"Участник завершил основной круг","time":"10:32:22.472","competitorID":5}
   {"time":"2025-05-07T03:19:51.3194668+03:00","level":"INFO","msg":"Generating final report"}
   {"time":"2025-05-07T03:19:51.3194668+03:00","level":"INFO","msg":"Application completed successfully"}
   Final Results:
   ID  Status    Total Time    Laps Times                  Speed Laps    Penalty Times               Speed Penalty  Hits/Shots
   --  ------    ----------    ----------                  ----------    -------------               -------------  ----------
   2   Finished  00:25:18.356  00:12:39.746, 00:12:38.610  4.607, 4.614  00:00:50.000, 00:00:50.000  3.000, 3.000   8/10
   . . 
   . . .
   ```

---
### Configuration

Sample _config.json_ in `input/config/`:
```
{
"laps": 2,
"lapLen": 3651,
"penaltyLen": 50,
"firingLines": 1,
"start": "09:30:00.000",
"startDelta": "00:00:30"
}
```

---
### Event File

Sample _events_ file in `input/events/events`:
```
[09:31:49.285] 1 3
[09:32:17.531] 1 2
[09:37:47.892] 1 5
[09:38:28.673] 1 1
[09:39:25.079] 1 4
[09:55:00.000] 2 1 10:00:00.000
[09:56:30.000] 2 2 10:01:30.000
. .
. . .
```
(see examples in README_TZ_*.md)

---
### Output Examples

1. Minimal report (`make run`):
   ```
   Final Results:
   [Finished] 2 [{10:14:09.746, 4.607}, {10:26:48.356, 4.614}] {00:01:40.000, 3.000} 8/10
   [Finished] 1 [{10:12:35.380, 4.633}, {10:25:26.047, 4.542}] {00:02:30.000, 3.000} 7/10
   [Finished] 3 [{10:15:43.273, 4.586}, {10:28:34.773, 4.537}] {-, -} 10/10
   [Finished] 4 [{10:17:16.947, 4.564}, {10:30:36.413, 4.378}] {00:01:40.000, 3.000} 8/10
   [Finished] 5 [{10:19:21.270, 4.368}, {10:32:22.472, 4.480}] {00:02:30.000, 3.000} 7/10
   ```

2. Full report (`make run-fullOutput`):
   ```
   Final Results:
   ID  Status    Total Time    Laps Times                  Speed Laps    Penalty Times               Speed Penalty  Hits/Shots
   --  ------    ----------    ----------                  ----------    -------------               -------------  ----------
   2   Finished  00:25:18.356  00:12:39.746, 00:12:38.610  4.607, 4.614  00:00:50.000, 00:00:50.000  3.000, 3.000   8/10
   1   Finished  00:25:26.047  00:12:35.380, 00:12:50.667  4.633, 4.542  00:01:40.000, 00:00:50.000  3.000, 3.000   7/10
   3   Finished  00:25:34.773  00:12:43.273, 00:12:51.500  4.586, 4.537                                             10/10
   4   Finished  00:26:06.413  00:12:46.947, 00:13:19.466  4.564, 4.378  00:01:40.000                3.000          8/10
   5   Finished  00:26:22.472  00:13:21.270, 00:13:01.202  4.368, 4.480  00:01:40.000, 00:00:50.000  3.000, 3.000   7/10
   ```
____

$${\color{green}> \space \color{green}Russian \space \color{green}version}$$
---:
<br>

Прототип системы для управления и анализа биатлонных соревнований.
Обрабатывает конфигурационные файлы и события, генерирует итоговые отчеты.<br>

**Описание задания:**
- [Описание задания на русском](README_TZ_rus.md)
- [Task description in English](README_TZ_eng.md)

---
## Установка
1. Убедитесь, что установлен Go версии 1.20 или новее.
2. Клонируйте репозиторий:
   - git clone https://github.com/BiathlonRaceProto-Yadro.git
   - cd BiathlonRaceProto-Yadro
3. Установите зависимости:
   make deps

---
## Запуск через Makefile

**Команды:**
- ``make build``        - Собрать проект
- ``make run``          - Запустить с минимальным выводом (как по заданию XD)
- ``make run-fullOutput`` - Запустить с полным табличным отчетом (работает по умолчанию при вызове ___make___)
- ``make clean``        - Удалить собранные файлы
- ``make setup``        - Подготовить каталоги и примеры конфигов

**Примеры:**
```
~\GolandProjects\BiathlonRaceProto-Yadro git:[main]
make run
"Сборка проекта..."
go build -o ./cmd/run/biathlon ./cmd/run/main.go
"Запуск приложения в minimal версии..."
./cmd/run/biathlon ./input/config/config.json ./input/events/events
Final Results:
. . .
```

---
## Ручной запуск

1. Сборка:
   - или это`go build -o <название_исполняемого_файла> cmd/run/main.go`
   - или можно скипнуть и перейти к _п.2_<br><br>

   Пример:
   ```
   ~\GolandProjects\BiathlonRaceProto-Yadro\cmd\run git:[main]
   go build -o biathlon main.go
   ```
2. Базовый запуск:
   - или это `./biathlon <путь_к_config.json> <путь_к_событиям>`
   - или это, если скинул _п.1_`./go run main.go <путь_к_config.json> <путь_к_событиям>`<br><br>

   Пример:
   ```
   ~\GolandProjects\BiathlonRaceProto-Yadro\cmd\run git:[main]
   go run main.go ..\..\input\config\config.json ..\..\input\events\events
   Final Results:
   . . .
   ```

3. Запуск с флагами:
   - `-debug`       - Включить отладочные логи
   - `-info`        - Логи уровня INFO
   - `-error`       - Логи уровня ERROR
   - `-fullOutput`  - Полный табличный отчет <br><br>

   Пример:
   ```
   ~\GolandProjects\BiathlonRaceProto-Yadro\cmd\run git:[main]
   go run main.go -fullOutput -info ..\..\input\config\config.json ..\..\input\events\events
   {"time":"2025-05-07T03:19:51.2920861+03:00","level":"INFO","msg":"Участник зарегистрирован","time":"09:31:49.285","competitorID":3}
   . . 
   . . .
   {"time":"2025-05-07T03:19:51.3194668+03:00","level":"INFO","msg":"Участник завершил основной круг","time":"10:32:22.472","competitorID":5}
   {"time":"2025-05-07T03:19:51.3194668+03:00","level":"INFO","msg":"Generating final report"}
   {"time":"2025-05-07T03:19:51.3194668+03:00","level":"INFO","msg":"Application completed successfully"}
   Final Results:
   ID  Status    Total Time    Laps Times                  Speed Laps    Penalty Times               Speed Penalty  Hits/Shots
   --  ------    ----------    ----------                  ----------    -------------               -------------  ----------
   2   Finished  00:25:18.356  00:12:39.746, 00:12:38.610  4.607, 4.614  00:00:50.000, 00:00:50.000  3.000, 3.000   8/10
   . . 
   . . .
   ```

---
### Конфигурация

Исходный _config.json_ в input/config/:
```
{
"laps": 2,
"lapLen": 3651,
"penaltyLen": 50,
"firingLines": 1,
"start": "09:30:00.000",
"startDelta": "00:00:30"
}
```

---
### Файл событий

Исходный файл _events_ в input/events/events:
```
[09:31:49.285] 1 3
[09:32:17.531] 1 2
[09:37:47.892] 1 5
[09:38:28.673] 1 1
[09:39:25.079] 1 4
[09:55:00.000] 2 1 10:00:00.000
[09:56:30.000] 2 2 10:01:30.000
. .
. . .
```
(см. примеры в README_TZ_*.md)

---
### Примеры вывода

1. Минимальный отчет (make run):
   ```
   Final Results:
   [Finished] 2 [{10:14:09.746, 4.607}, {10:26:48.356, 4.614}] {00:01:40.000, 3.000} 8/10
   [Finished] 1 [{10:12:35.380, 4.633}, {10:25:26.047, 4.542}] {00:02:30.000, 3.000} 7/10
   [Finished] 3 [{10:15:43.273, 4.586}, {10:28:34.773, 4.537}] {-, -} 10/10
   [Finished] 4 [{10:17:16.947, 4.564}, {10:30:36.413, 4.378}] {00:01:40.000, 3.000} 8/10
   [Finished] 5 [{10:19:21.270, 4.368}, {10:32:22.472, 4.480}] {00:02:30.000, 3.000} 7/10
   ```

2. Полный отчет (make run-fullOutput):
   ```
   Final Results:
   ID  Status    Total Time    Laps Times                  Speed Laps    Penalty Times               Speed Penalty  Hits/Shots
   --  ------    ----------    ----------                  ----------    -------------               -------------  ----------
   2   Finished  00:25:18.356  00:12:39.746, 00:12:38.610  4.607, 4.614  00:00:50.000, 00:00:50.000  3.000, 3.000   8/10
   1   Finished  00:25:26.047  00:12:35.380, 00:12:50.667  4.633, 4.542  00:01:40.000, 00:00:50.000  3.000, 3.000   7/10
   3   Finished  00:25:34.773  00:12:43.273, 00:12:51.500  4.586, 4.537                                             10/10
   4   Finished  00:26:06.413  00:12:46.947, 00:13:19.466  4.564, 4.378  00:01:40.000                3.000          8/10
   5   Finished  00:26:22.472  00:13:21.270, 00:13:01.202  4.368, 4.480  00:01:40.000, 00:00:50.000  3.000, 3.000   7/10
   ```



