# Распределенный вычислитель арифметических выражений

<details>
  <summary>Техническое задание</summary>
  
  Пользователь хочет считать арифметические выражения. Он вводит строку 2 + 2 * 2 и хочет получить в ответ 6. Но наши операции сложения и умножения (также деления и вычитания) выполняются "очень-очень" долго. Поэтому вариант, при котором пользователь делает http-запрос и получает в качетсве ответа результат, невозможна. Более того: вычисление каждой такой операции в нашей "альтернативной реальности" занимает "гигантские" вычислительные мощности. Соответственно, каждое действие мы должны уметь выполнять отдельно и масштабировать эту систему можем добавлением вычислительных мощностей в нашу систему в виде новых "машин". Поэтому пользователь, присылая выражение, получает в ответ идентификатор выражения и может с какой-то периодичностью уточнять у сервера "не посчиталость ли выражение"? Если выражение наконец будет вычислено - то он получит результат. Помните, что некоторые части арфиметического выражения можно вычислять параллельно.

Front-end часть

GUI, который можно представить как 4 страницы

Форма ввода арифметического выражения. Пользователь вводит арифметическое выражение и отправляет POST http-запрос с этим выражением на back-end. Примечание: Запросы должны быть идемпотентными. К запросам добавляется уникальный идентификатор. Если пользователь отправляет запрос с идентификатором, который уже отправлялся и был принят к обработке - ответ 200. Возможные варианты ответа:
200. Выражение успешно принято, распаршено и принято к обработке
400. Выражение невалидно
500. Что-то не так на back-end. В качестве ответа нужно возвращать id принятного к выполнению выражения.
Страница со списком выражений в виде списка с выражениями. Каждая запись на странице содержит статус, выражение, дату его создания и дату заверщения вычисления. Страница получает данные GET http-запрсом с back-end-а
Страница со списком операций в виде пар: имя операции + время его выполнения (доступное для редактирования поле). Как уже оговаривалось в условии задачи, наши операции выполняются "как будто бы очень долго". Страница получает данные GET http-запрсом с back-end-а. Пользователь может настроить время выполения операции и сохранить изменения.
Страница со списком вычислительных можностей. Страница получает данные GET http-запросом с сервера в виде пар: имя вычислительного ресурса + выполняемая на нём операция.

Требования:
Оркестратор может перезапускаться без потери состояния. Все выражения храним в СУБД.
Оркестратор должен отслеживать задачи, которые выполняются слишком долго (вычислитель тоже может уйти со связи) и делать их повторно доступными для вычислений.

Back-end часть

Состоит из 2 элементов:

Сервер, который принимает арифметическое выражение, переводит его в набор последовательных задач и обеспечивает порядок их выполнения. Далее будем называть его оркестратором.
Вычислитель, который может получить от оркестратора задачу, выполнить его и вернуть серверу результат. Далее будем называть его агентом.
Оркестратор
Сервер, который имеет следующие endpoint-ы:

Добавление вычисления арифметического выражения.
Получение списка выражений со статусами.
Получение значения выражения по его идентификатору.
Получение списка доступных операций со временем их выполения.
Получение задачи для выполения.
Приём результата обработки данных.

Агент
Демон, который получает выражение для вычисления с сервера, вычисляет его и отправляет на сервер результат выражения. При старте демон запускает несколько горутин, каждая из которых выступает в роли независимого вычислителя. Количество горутин регулируется переменной среды.
</details>


## Как запускать?
  К сожалению я так и не смог насторить докер для запуска, проблема возникает на этапе загрузки миграции для постгреса, так что проект должен был запускаться по команде ```docker compose up -d``` но увы, так не запустится.<br>
  ### Действия для запуска:
  1. Вам подадобится скачать <a href="https://docs.docker.com/get-docker/" >docker</a>
  2. Запустите docker
  3. Запустите команду: ```docker run -d --hostname my-rabbit --name some-rabbit -p 15672:15672 -p 5672:5672 rabbitmq:3-management``` для создания контейнера RabbitMQ
  4. Запустите команду: ```docker run --name my-postgres -p 5432:5432 -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -e POSTGRES_DB=MicroserviceCalculatorDB -d postgres``` для запуска postgresql
  5. Зайдите в docker desktop, во вкладку *containers* и убедитесь что контейнер *some-rabbit* и *my-postgres* запущен (если не запущены - запустите кнопкой)
  6. Откройте терминал в корневой дирриктории проекта (*MicroserviceCalculatorProject*)
  7. **Важно** чтоб на компьютере стоял go
  8. в терминале запустите команду ```go install github.com/pressly/goose/v3/cmd/goose@latest```
  9. Перейдите в терминале в *MicroserviceCalculatorProject/internal/sql/schema* и запустите команду ```goose postgres postgres://admin:admin@localhost:5432/MicroserviceCalculatorDB up``` (8 и 9 пункт нужны для настройки базы данных (создания таблиц))
  10. Запустите файлы *agent.exe* и *orchestrator.exe* которые лежат в *MicroserviceCalculatorProject/cmd*(агентов можно запустить несколько, это можно сделать через командную строку или через проводник запустить exe файлы).
  11. Откройте командную строку и введите команду ```curl localhost:8080/api``` и вам откроется минидокументация<br>

  Если по какой-то причине что-то не получилось, свяжитесь со мной.

  Если вы пользователь Линукс/МакОС то *.exe* файлы у вас скорее всего не запустятся. Вам прийдется скачать все пакеты которые прописанны в *go.mod* и в разных термаиналах запустить команды:
  * ```go run <путь до проекта>/MicroserviceCalculatorProject/orchestrator/cmd/cmd.go```
  * ```go run <путь до проекта>/MicroserviceCalculatorProject/agent/cmd/cmd.go```


## Как устроен проект
  В проеке есть две основные папки: orchestrator и agent. В каждой из них находятся папки cmd (вход в программу), internal (внутренние методы/структуры) и pkg (файлы которые могут пригодится в любом проекте). в orchestrator/sql хранятся миграции (схемы) и запросы для работы с postgreSQL.

## Технологии
  в проекте используются:
  * *RabbitMQ* для "общения" оркестратора и агентов
  * *PostgreSQL* для хранения результатов выражений/подвыражений и долгота арифметических действий
  * *chi* пакет для обработки http-запросов (роутер)

## Конвертация выражения
  1. Выражение приводится с стандартному виду: убираются лишние пробелы, разбивается посимвольно и слайс символов передается на пункт *2* (ф-ия FormatExpression).
  2. Происходит проверка валидации и если все хорошо проходит дальше, иначе возвращает пользователю ошибку (ф-ия IsValid).
  3. Выражение переводится в постфиксную запись (ф-ия infixExpToPostfixExp).
  4. Из постфиксной записи происходит конвертация в дерево AST (ф-ия buildAST) 
  5. Дерево AST рекурсивно группируется в мапу подвыражений, где ключи - номер подвыражения, а значения - сами подвыражения. Разбивается на элементарные подвыражения (два операнда и оператор) а неизвестные операнды заменяются номерами подвыражеий, в которых они подсчитаются (splitIntoSubexpressions).

    *Пример:*
      2 +2*  2 ->
      2 + 2 * 2 ->
      ["2", "+", "2", "*", "2"] ->
      ["2", "2", "*", "2", "+"] ->
      treeAST (не отрисовать) ->
      map{
        1: "2 * 2",
        2: "{1} + 2"
      }


    

## Некоторые примеры и запросы
 * curl -X POST -H "Content-Type: application/json" -d "{\\"expression\\":\\"(4 + 2) + 5 * 6\\"}" localhost:8080/api/expression
  * curl -X POST -H "Content-Type: application/json" -d "{\\"expression\\":\\"(2 + 2 + 2 + 2) / 4\\"}" localhost:8080/api/expression
  * curl -X POST -H "Content-Type: application/json" -d "{\\"expression\\":\\"(2 + 2) * 4 + 3 - 4 + 5\\"}" localhost:8080/api/expression 

  #### Изменение продолжительности операций:
  Для этого требуется отправить POST запрос по адресу localhost:8080/api/durations с указанием имя оператора и продолжительность вычисления.<br>
  *Пример:*

    curl -X POST -H "Content-Type: application/json" -d "{\"operator\": \"*\",\"duration\": 5}" localhost:8080/api/durations

## Схемы
![schema](images/schema.png)
![cshema](images/orchestrator.png)

Ссылка на схемы: https://excalidraw.com/#json=nmhdkFY2NMesc0_JL6Eag,JYf3s11swuItgDbwwoA5rg

-----
Простите, многое не успел написать, по всем вопросам пишите мне в telegram: https://t.me/Ruslan20007

