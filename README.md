<h1 align="center">
  <img height="180px" src="https://github.com/raito-io/docs/raw/main/images/Raito_Logo_Vertical_RGB.png" alt="Raito" />
</h1>

<h4 align="center">
  Extensible CLI to easily manage the access controls for your data sources.
</h4>

<p align="center">
    <a href="/LICENSE.md" target="_blank"><img src="https://img.shields.io/badge/license-Apache%202-brightgreen.svg" alt="Software License" /></a>
    <a href="https://github.com/raito-io/cli/actions/workflows/build.yml" target="_blank"><img src="https://img.shields.io/github/workflow/status/raito-io/cli/Raito%20CLI%20-%20Build/main" alt="Build status" /></a>
    <a href="https://codecov.io/gh/raito-io/cli" target="_blank"><img src="https://img.shields.io/codecov/c/github/raito-io/cli" alt="Code Coverage" /></a>
</p>

<hr/>

:rotating_light: :rotating_light: :rotating_light:  

**Note: This repository is still in a very early stage of development.  
It contains code that will allow communication with Raito Cloud once it is released. 
At this point, no contributions are accepted to the project yet.**  

:rotating_light: :rotating_light: :rotating_light:

# Introduction
This is the core Raito CLI implementation.

# Install
Using HomeBrew:
```bash
brew install raito-io/tap/cli
```

# Usage
To get an overview of the possibilities, simply use the `--help` flag
```bash
raito --help
```

<!--
# Join the Raito Community
We would love to hear your thoughts and questions.  
So please join our [Slack Community](https://join.slack.com/t/raitocommunity/shared_invite/zt-13ti14ezm-RsGFyJq4FU9IEfjqg_POag) if you would like to join the conversation and contribute.
-->
# Contributing
<!--
Want to contribute to the Raito open source code base? Great!  
Take a look at our [Contribution Guide](https://github.com/raito-io/cli/blob/HEAD/CONTRIBUTING.md) to get you started.
-->
### Prerequisites
 - Install the correct version of Go (see go.mod for the version being used now)
 - After you check out this git repository, make sure to execute the following command (in the root of the repository) to run the pre-commit hooks: `git config core.hooksPath .githooks`
 - We use `golangci-lint` to check the code for quality. Please make sure to [install it](https://golangci-lint.run/usage/install/#local-installation)
