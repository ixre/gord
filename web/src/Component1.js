import React from "react"
import Icon from "antd/lib/icon"
class Component1 extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        return (
            <div>
                <Icon type="home"/>
                <div class="c1">Hello {this.props.name}</div>
            </div>
        )
    };
}
export default Component1;
