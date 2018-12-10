import React from 'react'
import DisplayDate from './display_date.js'
import './date_range.css'

const DateRange = ({startDate, endDate}) => {
  let {dateTime: startDateTime, resolution: startResolution} = startDate
  let {dateTime: endDateTime, resolution: endResolution} = endDate
  return (
    <span className="date-range">
      <DisplayDate dateTime={startDateTime} resolution={startResolution}/><span className="date-range-separator">-</span><DisplayDate dateTime={endDateTime} resolution={endResolution}/>
    </span>
  )
}

export default DateRange
