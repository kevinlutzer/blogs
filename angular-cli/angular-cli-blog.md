# Integrating the Angular Cli with the VA Tech Stack

## Basics

Install the angular-cli by running
```bash
npm install -g @angular/cli
```

## Integration With Vendasta's FE Tech Stack

A generated angular-cli app does not contain any of our usual front-end continuous integration configuration. Nor does it contain the basic test setup we use

## Using Yarn

Globally configure your version of the angular-cli to use yarn as the default package manager, by running `ng set --global packageManager=yarn` in a terminal window. To set yarn as the default package manager -- to be used during contiunous integration -- for an angular-cli app, run `ng set packageManager=yarn`. This command sets the `packageManager` key in `.angular-cli.json` to be `yarn`. If you run this command after installing the `node_modules` with npm, you will have to delete them.

`ng set --global packageManager=yarn`
`ng set packageManager=yarn`

## Spec (Jasmine) Tests
The base angular-cli app is configured to work with [karma](https://karma-runner.github.io/1.0/index.html), using the Chrome launcher to run the server. This works well for local development, however our Jenkins CI container is built to work with the [phantomjs](https://www.npmjs.com/package/karma-phantomjs-launcher) launcher. The following steps are to setup an angular-cli project with phatomjs.

### Setup

#### Install Depenencies

- Add phatomjs dependencies.
- Remove uneeded chrome launcher.

`yarn add intl` <br>
`yarn add karma-phantomjs-launcher --dev` <br>
`yarn remove karma-chrome-launcher --dev`

#### Edit karma.conf.json

- Replace the requirement `karma-chrome-launcher` with `karma-phantomjs-launcher`.
- Replace the browser, that hosts the karma server, `Chrome` with `PhantomJS`.

```json
module.exports = function (config) {
  config.set({
    ...
    plugins: [
      ...
      require('karma-phantomjs-launcher')
    ],
    ...
    browsers: ['PhantomJS'],
    ...
  });
};
```

#### Edit polyfill.js

- Add an import for the [intl](https://developer.mozilla.org/en/docs/Web/JavaScript/Reference/Global_Objects/Intl) package. This package is needed as phatomjs does not implement the ECMAScript Internationalization API which is required for basic JS operations.
- Add an import for the core-js `shim`. 

`import 'core-js/client/shim';` <br>
`import 'intl';`

#### Test the Testing Setup

`yarn test --browsers PhantomJS --singleRun true`

## E2E tests

By default the angular-cli uses [protractor](http://www.protractortest.org/#/) to run 