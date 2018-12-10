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
    image {
      thumbnail
    }
    affiliationList {
      id
      label
    }
  }
}
`
let allPeople = () => {
  return apollo.query({
    query: ALL_PEOPLE
  })
}

const PERSON_BY_ID = gql`
query personById($id: String) {
  person(id: $id) {
    id
    name {
      firstName
      lastName
      middleName
    }
    image {
      main
    }
    affiliationList {
      id
      label
      startDate {
        dateTime
        resolution
      }
    }
    overviewList {
      overview
      type {
        code
        label
      }
    }
    educationList {
      label
      org {
        id
        label
      }
    }
    publicationList {
      id
      label
      venue {
        label
        uri
      }
      authorList
    }
  }
}
`

let personById = (id) => {
  return apollo.query({
    query: PERSON_BY_ID,
    variables: {
      id: id
    }
  })
}

export default {
  allPeople,
  personById
}


