import React from 'react'
import { format } from 'date-fns'
import './display_date.css'

const DisplayDate = ({dateTime, resolution}) => {
  if (!dateTime) return null
  let d = new Date(dateTime)
  let displayDate
  if (resolution === 'year') {
    displayDate = format(d,'YYYY')
  } else if (resolution === 'yearMonth') {
    displayDate = format(d,'MMM YYYY')
  } else if (resolution === 'yearMonthDay') {
    displayDate = format(d,'MMM D, YYYY')
  } else {
    displayDate = `unknown date resolution: ${resolution}`
  }
  return (
    <span className="display-date">
      {displayDate}
    </span>
  )
}

export default DisplayDate
