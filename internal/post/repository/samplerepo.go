package repository

import (
	"sort"

	"github.com/2024_2_BetterCallFirewall/internal/post/models"
)

type Repository struct {
	storage []*models.Post
}

func NewRepository() *Repository {
	return &Repository{
		storage: []*models.Post{
			{Header: "Alexey Zemliakov", Body: "Go (часто также golang) — компилируемый многопоточный язык программирования, разработанный внутри компании Google. Поддерживает объектно-ориентированный, и функциональный стили. Разработка Go началась в сентябре 2007 года, его непосредственным проектированием занимались Роберт Гризмер, Роб Пайк и Кен Томпсон, занимавшиеся до этого проектом разработки операционной системы Inferno. Официально язык был представлен в ноябре 2009 года", CreatedAt: "2024-09-20"},
			{Header: "Andrew Savvateev", Body: "В России альпинизм, спортивная сущность которого состоит в преодолении естественных препятствий при восхождении на вершины гор, является официально признанным видом спорта и активного отдыха. Как в советский так и постсоветский периоды отношение к альпинизму было как к одному из видов спорта, имеющего большое прикладное значение. ", CreatedAt: "2024-09-21"},
			{Header: "Konstantin Galanin", Body: "Figma — онлайн-сервис для разработки интерфейсов и прототипирования с возможностью организации совместной работы в режиме реального времени. Используется как для создания упрощённых прототипов интерфейсов, так и для детальной проработки дизайна интерфейсов мобильных приложений, веб-сайтов, корпоративных порталов.\n\nСервис доступен по подписке, предусмотрен бесплатный тарифный план для одного пользователя. Имеются небраузерные онлайн-версии для Windows, macOS. Реализована интеграция с корпоративным мессенджером Slack и инструментом прототипирования Framer", CreatedAt: "2024-09-22"},
			{Header: "Alexander Novikov", Body: "JavaScript (англ. /ˈdʒɑːvəskrɪpt/; аббр. JS) — мультипарадигменный язык программирования. Поддерживает объектно-ориентированный, императивный и функциональный стили. Является реализацией спецификации ECMAScript (стандарт ECMA-262).\n\nJavaScript обычно используется как встраиваемый язык для программного доступа к объектам приложений. Наиболее широкое применение находит в браузерах как язык сценариев для придания интерактивности веб-страницам.\n\nОсновные архитектурные черты: динамическая типизация, слабая типизация, автоматическое управление памятью, прототипное программирование, функции как объекты первого класса. ", CreatedAt: "2024-09-23"},
			{Header: "Алексей Земляков", Body: "Мощное гравитационное поле замедляет время\nИз-за гравитации время в космосе протекает по-разному. Чем мощнее гравитационное поле, тем сильнее замедляется время. Этот феномен проиллюстрирован в фильме «Интерстеллар» Кристофера Нолана. Когда герои попадают на планету Миллер, час для них оказывается равен семи земным годам. Вернувшись на борт космического корабля спустя три с небольшим часа, астронавты застают уже поседевшего коллегу, который ждал их возвращения долгие 23 года. Практически так же происходит и в реальности. Например, для космонавтов время тянется на доли секунды быстрее, чем для людей на Земле. А вблизи черной дыры оно почти полностью останавливается.", CreatedAt: "2024-09-24"},
			{Header: "Андрей Савватеев", Body: "В одном из исследований было обнаружено, что альпинисты тратят больше времени на выбор одежды, чем на само восхождение. Ведь как известно, какой цвет куртки и штанов лучше сочетается с горным пейзажем - это важный фактор успеха!\n\nНекоторые альпинисты настолько привязаны к своим веревкам, что никогда не разлучаются с ними даже в повседневной жизни. Они идут с веревкой по улице, в магазин, и даже на свидание!\n\nОказывается, альпинисты не только мастера восхождений, но и настоящие профессионалы в игре \"Угадай, где я\". Они могут увидеть фотографию горы и сразу определить, на какой высоте они находятся, какой маршрут использовался и какие препятствия им пришлось преодолеть.", CreatedAt: "2024-09-25"},
			{Header: "Константин Галанин", Body: "Вы видели новые кампусы?\nОни просто шедевральны, современный дизайн, благоустройство, все на высшем уровне!", CreatedAt: "2024-09-26"},
			{Header: "Новиков Александр", Body: "Факультет СМ - лучший факультет в МГТУ им. Н.Э.Баумана.\n Поступая сюда вы точно получите отличное образование в области инженерии. Этот факультет содержит огромное количество направлений обучения на которых обучают всему, начиная от построения ракет для покорения космоса, вплоть до подводных и наземных роботов.", CreatedAt: "2024-09-27"},
			{Header: "Анекдоты", Body: "Царь позвал к себе Иванушку-дурака и говорит:\n— Если завтра не принесёшь двух говорящих птиц, голову срублю.\nИван принёс филина и воробья. Царь велит:\n— Ну, пусть что-нибудь скажут.\nИван спрашивает:\n— Воробей, почем раньше водка в магазине была?\n— Чирик.\n— Филин, подтверди.\n— Подтверждаю.", CreatedAt: "2024-09-25"},
			{Header: "Технопарк", Body: "Всем привет!\nПодходит к концу первый модуль, грядёт рубежный контроль, а мы должны применить на практике знания, полученные на предыдущей лекции.\n \nНеобходимо сделать следующее:\n\n    Перепишите реализацию модуля, который упрощает работу с HTTP-запросами и необходимые для работы приложения сервисы, с использованием промисов. За использование Fetch API и async/await будет респект от преподавателей. \n    Договоритесь в команде об API и задокументируйте его, используя любой инструмент. Почитайте о формате комментариев в JavaScript-коде JSDoc, документируйте ваши классы и модули с его использованием. Использование инструментов для генерации документации не требуется, но поощряется доп. баллами :]\n    Реализуйте интеграцию между фронтендом и бекендом. Настройте CORS\n    За доп. баллы вы можете внедрить в ваши приложения защиту от CSRF-атак\n", CreatedAt: "2024-09-29"},
		},
	}
}

func (r *Repository) GetAll() []*models.Post {
	sort.Slice(r.storage, func(i, j int) bool {
		return r.storage[i].CreatedAt > r.storage[j].CreatedAt
	})
	return r.storage
}