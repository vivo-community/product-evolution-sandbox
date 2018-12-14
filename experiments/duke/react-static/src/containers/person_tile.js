import React from 'react'
import { Link } from '@reach/router'
import './person_tile.css'

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

  let givenName = [firstName,middleName].filter((n) => n).join(" ")
  let displayName = `${lastName}, ${givenName}`
  let displayTitle = affiliationList[0].label
  return (
    <div className="person-tile">
      <Link to={`/people/${id}`}>
        <img className="person-tile-image" src={thumbnail} alt={`${displayName} Thumbnail`}/>
        <div className="person-tile-name">{displayName}</div>
        <div className="person-tile-title">{displayTitle}</div>
      </Link>
    </div>
  )

}

export default PersonTile
