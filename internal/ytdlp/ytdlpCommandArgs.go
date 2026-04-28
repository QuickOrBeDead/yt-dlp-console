package ytdlp

type YtDlpCommandArgs struct {
	url      string
	password string
	args     []string
}

func NewYtDlpCommandArgs(url, password string) *YtDlpCommandArgs {
	return &YtDlpCommandArgs{
		url:      url,
		password: password,
		args:     make([]string, 0),
	}
}

func (c *YtDlpCommandArgs) AddArg(arg string) {
	c.args = append(c.args, arg)
}

func (c *YtDlpCommandArgs) AddArgWithValue(arg, value string) {
	c.args = append(c.args, arg)
	c.args = append(c.args, value)
}

func (c *YtDlpCommandArgs) BuildArgs() []string {
	return c.buildArgs(true)
}

func (c *YtDlpCommandArgs) BuildArgsMasked() []string {
	return c.buildArgs(false)
}

func (c *YtDlpCommandArgs) ClearPassword() {
	if len(c.password) > 0 {
		c.password = ""
	}
}

func (c *YtDlpCommandArgs) buildArgs(showPassword bool) []string {
	args := make([]string, 0)

	if len(c.password) > 0 {
		if showPassword {
			args = append(args, "--video-password", c.password)
		} else {
			args = append(args, "--video-password", "******")
		}
	}

	args = append(args, c.args...)
	args = append(args, c.url)

	return args
}
