import axios from 'axios'
import Scholars from './src/fetchers/scholars'
import { makePageRoutes } from 'react-static/node'

export default {
  getSiteData: () => ({
    title: 'Scholars Graphql Example',
  }),
  getRoutes: async () => {
    const { data: people } = await Scholars.allPeople()
    const allPeople = people.personList.filter((p) => p.image.thumbnail != "")
    return [
      {
        path: '/',
        component: 'src/containers/Home',
      },
      ...makePageRoutes({
        items: allPeople,
        pageSize: 40,
        pageToken: 'page',
        route: {
          path: '/people',
          component: 'src/containers/people',
        },
        decorate: (items, i, totalPages) => ({
          getData: () => ({
            people: items,
            currentPage: i,
            totalPages
          })
        })
      }),
      {
        is404: true,
        component: 'src/containers/404',
      },
      {
        path: '/people',
        children: allPeople.map(person => ({
          path: `/${person.id}`,
          component: 'src/containers/person',
          getData: async () => {
            let {data} = await Scholars.personById(person.id)
            return data
          }
        }))
      }
    ]
  },
  // webpack: (config, { defaultLoaders }) => {
  //   config.module.rules = [
  //     {
  //       oneOf: [
  //         {
  //           test: /\.json$/,
  //           use: [{ loader: 'json-loader' }],
  //         },
  //         defaultLoaders.jsLoader,
  //         defaultLoaders.cssLoader,
  //         defaultLoaders.fileLoader,
  //       ],
  //     },
  //   ]
  //   return config
  // },
}
