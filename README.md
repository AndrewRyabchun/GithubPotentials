# Github Potentials

Github Potentials command line tool is package that allows you to find rising Github repositories in 3 steps:
  - Fetch recently updated repositories (up to 1k).
  - Count stats for last *n* hours/days: new stars, unique contributors and commits.
  - Sort by selected criteria and take the best.

### Command line tool
Under [cmd/ghp](https://github.com/ArtIsResistance/GithubPotentials/tree/master/cmd/ghp) directory you can see an example of using this package.

`go get github.com/artisresistance/githubpotentials/cmd/ghp`

You must provide config file that contains your Github API secret token.
Example output you can find [here](https://githubpotentials.azurewebsites.net/data.json).