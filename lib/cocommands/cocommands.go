package cocommands

func info() string {
  infoString := `


                             ______   ____     ____     ____   _____
                            / ____/  / __ \   / __ \   /  _/  / ___/
                           / /      / / / /  / / / /   / /    \__ \ 
                          / /___   / /_/ /  / /_/ /  _/ /    ___/ / 
                          \____/   \____/  /_____/  /___/   /____/  
             --------------------------------------------------------------------
             Codis is an open-source (GPL-v3) search engine made by Timo Kats. 
             This version is meant as a proof-of-concept. As a result, feedback
             is highly appreciated. So, if you have any, please send an e-mail 
             to tpakats@gmail.com or visit my GitHub at github.com/TimoKats. Thx!

             codis v0.0.1 January 2024
  `
  return infoString
}

func ParseCommand(command string) ([]string, []string) {
  if command == "info" {
    return []string{info()}, []string{"info"}
  } else if command == "help" {
    return []string{"HELP PAGE!"}, []string{"help"}
  } else {
    return []string{"invalid command"}, []string{"invalid command"}
  }
}
