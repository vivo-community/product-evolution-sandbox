import React from 'react'

import { render, unmountComponentAtNode  } from 'react-dom'
import VivoSearch from './containers/VivoSearch'

require ('bootstrap')

import 'jquery'
import "babel-polyfill"

// NOTE: wanted to require this in particular tabs (where actually needed)
// but babel-node tries to parse *.less as *.js file
require('./styles/vivo_search.scss');

module.exports = function(targetNode) {
  unmountComponentAtNode(targetNode)
  render (
      <VivoSearch />,
      targetNode
  )
}


