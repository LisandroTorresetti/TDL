package bot

type summarizer interface {
	SummarizeNews(news string) (string, error)
}

type NewsBot struct {
	gpt summarizer
}

func newNewsBot(gpt summarizer) *NewsBot {
	return &NewsBot{
		gpt: gpt,
	}
}

func (nb *NewsBot) Summarize(newsToSummarize string) (string, error) {
	return nb.gpt.SummarizeNews(newsToSummarize)
}
