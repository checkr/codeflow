# Codeflow Web Dashboard

This project template was built with [Create React App](https://github.com/facebookincubator/create-react-app).

## Available Scripts

In the project directory, you can run:

### `npm start`

Runs the app in the development mode.<br>
Open [http://localhost:3000](http://localhost:3000) to view it in the browser.

The page will reload if you make edits.<br>
You will also see any lint errors in the console.

### `npm run build`

Builds the app for production to the `build` folder.<br>
It correctly bundles React in production mode and optimizes the build for the best performance.

The build is minified and the filenames include the hashes.<br>
Your app is ready to be deployed!

### Updating npm packages 
```
$ npm outdated
$ npm outdated --depth=0 | grep -v Package | awk '{print $1}' | xargs -I% npm install %@latest --save
```
or [npm-check-updates](https://www.npmjs.com/package/npm-check-updates)

