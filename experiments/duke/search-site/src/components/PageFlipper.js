import React, { Component } from 'react'
import PropTypes from 'prop-types'
import ReactDOM from 'react-dom'

export default class PageFlipper extends Component {

  constructor(props, context) {
    super(props, context)
    this.onClick = this.props.onClick
  }
  
  flip(pages, direction) {
    console.log(`trying to make prev/next ${pages}:${direction}`)
     if(pages[0] == '+') {
        let pageNumber = x[1]

        let desc = (<span><span aria-hidden="true">&laquo;</span> Previous</span>)
        if (direction == 'forward') {
          desc = (<span>Next <span aria-hidden="true">&raquo;</span></span>)
        }
        let key = `pageLinkTo_${pageNumber}`
        return (
          <li key={key} className="page-info">
            <a href="#" className="page-link" onClick={(e) => this.onClick(e, pageNumber)}>{desc}</a>
          </li>
        ) 
      } else {
        return null
      }
  }

  render() {
    return this.flip(this.props.pages, this.props.direction)
  }
}


