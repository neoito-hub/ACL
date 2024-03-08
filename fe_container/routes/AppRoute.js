import React, { Suspense, lazy } from 'react'
import { Switch, Route } from 'react-router-dom'

const AppRoute = () => {
  const Home = lazy(() => import('../src/components/Home/home'))
  const UserProfile = lazy(() => import('../src/components/UserProfile'))
  const Invitation = lazy(() => import('../src/components/Invitation'))
  const CreateSpace = lazy(() => import('../src/components/CreateSpace'))
  const SpaceDetails = lazy(() => import('../src/components/SpaceDetails'))
  const SpaceListing = lazy(() => import('../src/components/SpaceListing'))

  return (
    <Suspense fallback="">
      <Switch>
        <Route exact path="/" component={Home} />
        <Route path="/home" component={Home} />
        <Route path="/profile" component={UserProfile} />
        <Route path="/invitation" component={Invitation} />
        <Route path="/spaces/:sub_menu" component={SpaceDetails} />
        <Route path="/spaces" component={SpaceListing} />
        <Route path="/create-space" component={CreateSpace} />
      </Switch>
    </Suspense>
  )
}

export default AppRoute
