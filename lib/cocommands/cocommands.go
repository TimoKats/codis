/* 
** @name: cocommands
** @author: Timo Kats
** @description: All the possible results from command mode 
*/

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

func help() []string {
  helpString := []string{
    `
    BASIC COMMANDS:
      <tab> to switch between query types
      / to enter query mode
      : to enter command mode
      <ctrl+f> to change settings
      <enter> to submit query
    `,
    `
    QUICK SEARCH:
      DESCRIPTION:
        Simply find the submitted text in the parent directories.
      COMMANDS:
        <enter> to search.
        <ctrl+k> or <ctrl+j> to iterate between results.
    `,
    `
    FUZZY SEARCH:
      DESCRIPTION:
        Can find instancies of an estimated/misspelled query.
      COMMANDS:
        <enter> to search.
        <ctrl+k> or <ctrl+j> to iterate between results.
    `,
    `
    EXPLORATIVE SEARCH:
      DESCRIPTION:
        Shows the filetree with info per file.
      COMMANDS:
        <enter>+integer to zoom in.
        <ctrl+d> to view directories only.
        <ctrl+g> to change displayed info.
        <ctrl+k> or <ctrl+j> to iterate between results.
    `,
    `
    DEPENDENCY SEARCH:
      DESCRIPTION:
        Shows the dependencies between the files using an indented filetree.
      COMMANDS:
        <ctrl+g> to change displayed info.
        <ctrl+k> or <ctrl+j> to iterate between results.
    `,
    `
    FILE VIEW:
      DESCRIPTION: 
        Displays a file based on a query (all are shown when query is empty) 
      COMMANDS:
        <ctrl+k> or <ctrl+j> to iterate between results.
        <ctrl+g> to change type of view.
    `,
  }
  return helpString
}

func ParseCommand(command string) ([]string, []string) {
  if command == "info" {
    return []string{info()}, []string{"info page"}
  } else if command == "help" {
    return help(), []string{"help page","help page","help page","help page","help page","help page",}
  } else {
    return []string{"invalid command"}, []string{"invalid command"}
  }
}
