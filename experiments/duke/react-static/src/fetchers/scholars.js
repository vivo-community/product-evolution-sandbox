import apollo from '../connectors/apollo'
import gql from 'graphql-tag'

const ALL_PEOPLE = gql`
query {
  personList {
    id
    name {
      firstName
      lastName
      middleName
    }
    affiliationList {
      id
      label
    }
    image {
      thumbnail
    }
  }
}
`
let allPeople = () => {
  return apollo.query({
    query: ALL_PEOPLE
  })
}

export default {
  allPeople
}


