import React from "react";
import Component1 from "./src/Component1";

class Login extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        return (
            <div>
                <Component1 name={"jarrysix"}/>
                <br/>
                <a class="a-about" onClick={() => {
                    this.props.history.push("/about")
                }}>关于</a>
            </div>
        );
    };
}

export default Login
