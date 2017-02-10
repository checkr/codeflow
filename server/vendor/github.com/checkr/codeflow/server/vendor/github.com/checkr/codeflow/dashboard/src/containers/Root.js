import '../../node_modules/bootstrap/dist/css/bootstrap.css'
import '../../public/css/application.css'
import '../../public/css/font-awesome.css'

if (process.env.NODE_ENV === 'production') {
  module.exports = require('./Root.prod')
} else {
  module.exports = require('./Root.dev')
}
