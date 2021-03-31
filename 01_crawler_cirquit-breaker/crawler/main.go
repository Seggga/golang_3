package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	// максимально допустимое число ошибок при парсинге
	errorsLimit = 100000

	// число результатов, которые хотим получить
	resultsLimit = 10000
)

var (
	// адрес в интернете (например, https://en.wikipedia.org/wiki/Lionel_Messi)
	url string

	// насколько глубоко нам надо смотреть (например, 3) и таймаут на работу парсера
	depthLimit, timeOut int
)

// Как вы помните, функция инициализации стартует первой
func init() {
	// задаём и парсим флаги
	flag.StringVar(&url, "url", "", "url address, no default value")
	flag.IntVar(&depthLimit, "depth", 3, "max depth for run, 3 by default")
	flag.IntVar(&timeOut, "timeout", 20, "timeout for run in seconds, 20 by default")
	flag.Parse()

	// Проверяем обязательное условие
	if url == "" {
		log.Print("no url set by flag")
		flag.PrintDefaults()
		os.Exit(1)
	}

	//url = "https://wikipedia.org/"
}

func main() {
	started := time.Now()

	ctx, cancel := context.WithCancel(context.Background())
	crawler := newCrawler(depthLimit)

	go watchSignals(cancel, crawler, timeOut)
	defer cancel()

	// создаём канал для результатов
	results := make(chan crawlResult)

	// запускаем горутину для чтения из каналов
	done := watchCrawler(ctx, results, errorsLimit, resultsLimit)

	// запуск основной логики
	// внутри есть рекурсивные запуски анализа в других горутинах
	crawler.run(ctx, url, results, 0)

	// ждём завершения работы чтения в своей горутине
	<-done

	log.Println(time.Since(started))
}

// ловим сигналы выключения и увеличения глубины поиска ссылок
func watchSignals(cancel context.CancelFunc, crawler *crawler, timeOut int) {
	// a channel to transmit quit-signals
	osSignalChan := make(chan os.Signal)
	// a channel to transmit "increase maxDepth"
	depthChan := make(chan os.Signal)

	signal.Notify(osSignalChan,
		syscall.SIGINT,
		syscall.SIGTERM)

	signal.Notify(depthChan, syscall.SIGHUP)

	timer := time.NewTimer(time.Second * time.Duration(timeOut))

	for {

		select {
		case <-timer.C:
			log.Printf("Parcer timeout (%d seconds)!", timeOut)
			cancel()
			return
		case sig := <-osSignalChan:
			log.Printf("got signal %q", sig.String())
			cancel()
			return
		case sig := <-depthChan:
			crawler.maxDepth += 2
			log.Printf("got signal %q, maximum depth has become %d", sig.String(), crawler.maxDepth)
		}
	}
}

func watchCrawler(ctx context.Context, results <-chan crawlResult, maxErrors, maxResults int) chan struct{} {
	readersDone := make(chan struct{})

	go func() {
		defer close(readersDone)
		for {
			select {
			case <-ctx.Done():
				return

			case result := <-results:
				if result.err != nil {
					maxErrors--
					if maxErrors <= 0 {
						log.Println("max errors exceeded")
						return
					}
					continue
				}

				log.Printf("crawling result: %v", result.msg)
				maxResults--
				if maxResults <= 0 {
					log.Println("got max results")
					return
				}
			}
		}
	}()

	return readersDone
}
