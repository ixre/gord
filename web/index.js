import "./index.css"
import Login from "./login"
import About from "./about"
import {HashRouter, Route, Switch} from 'react-router-dom'
import {Provider} from 'react-redux';

(async () => {
    let React = await import(/* webpackChunkName:"react" */ "react")
    let ReactDOM = await import(/* webpackChunkName:"react" */ "react-dom")

    let root = document.getElementById("root");
    ReactDOM.render((
        <HashRouter>
            <Switch>
                <Route exact path='/' component={Login}/>
                <Route path='/about' component={About}/>
            </Switch>
        </HashRouter>), root);
})();



