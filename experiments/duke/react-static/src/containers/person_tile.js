import React from 'react'
import { withRouteData, Link } from 'react-static'

const PersonTile = ({person}) => {
  let {
    id,
    name: {
      firstName,
      middleName,
      lastName
    },
    image: {
      thumbnail
    },
    affiliationList
  } = person

  let displayName = `${lastName}, ${firstName} ${middleName}`
  let displayTitle = affiliationList[0].label
  return (
    <div className="person-tile">
      <p>{displayName}</p>
      <p>{displayTitle}</p>
    </div>
  )

}

export default PersonTile
