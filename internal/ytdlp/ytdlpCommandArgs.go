package ytdlp

type YtDlpCommandArgs struct {
	url             string
	videoPassword   string
	username        string
	accountPassword string
	args            []string
}

func NewYtDlpCommandArgs(url, videoPassword string) *YtDlpCommandArgs {
	return &YtDlpCommandArgs{
		url:           url,
		videoPassword: videoPassword,
		args:          make([]string, 0),
	}
}

func (c *YtDlpCommandArgs) SetAccountAuth(username, accountPassword string) {
	c.username = username
	c.accountPassword = accountPassword
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
	c.videoPassword = ""
	c.accountPassword = ""
	c.username = ""
}

func (c *YtDlpCommandArgs) buildArgs(showPassword bool) []string {
	args := make([]string, 0)

	if len(c.username) > 0 {
		args = append(args, "--username", c.username)
		if showPassword {
			args = append(args, "--password", c.accountPassword)
		} else {
			args = append(args, "--password", "******")
		}
	}

	if len(c.videoPassword) > 0 {
		if showPassword {
			args = append(args, "--video-password", c.videoPassword)
		} else {
			args = append(args, "--video-password", "******")
		}
	}

	args = append(args, c.args...)
	args = append(args, c.url)

	return args
}
