var socket;

var Game = React.createClass({
   getInitialState: function () {
      return {
         players: null,
         planets: null,
         selectedPlanet: null
      };
   },

   unselectPlanet: function () {
      this.setState({ selectedPlanet: null });
   },

   setSelectedPlanet: function (planet) {
      this.setState({ selectedPlanet: planet });
   },

   getData: function () {
      var self = this;

      $.ajax({
         type: 'GET',
         url: '/data/',
         dataType: 'json',
         success: function(answer){
            var selectedPlanet = null
            if (self.state.selectedPlanet != null) { // gotta find this planet in the new planet list
               var selectedPlanetId = self.state.selectedPlanet.id;

               for (var i = 0; i < answer.planets.length; i++) {
                  if (answer.planets[i].id == selectedPlanetId) {
                     selectedPlanet = answer.planets[i];
                     break;
                  }
               }
            }

            // this is wrong: the references to the players and planets always change, so react keeps re-rendering it ARGXYZ
            self.setState({
               players: answer.players,
               planets: answer.planets,
               selectedPlanet: selectedPlanet
            });
         },
         error: function(xhr, type){
            console.log('Ajax error: GET /data/');
         }
      });
   },

   render: function () {
      return (
         <div>
            <StarMap planets={this.state.planets} players={this.state.players} setSelectedPlanet={this.setSelectedPlanet} />
            <ManagementInterface planets={this.state.planets} players={this.state.players} selectedPlanet={this.state.selectedPlanet} unselectPlanet={this.unselectPlanet} />
         </div>
      );
   },

   componentDidMount: function () {
      var self = this;
      socket = new WebSocket("ws://192.168.0.2:8081/ws/");

      socket.onopen = function (event) {
         if (document.cookie == null || document.cookie == 0) {
            var player = {
               name: prompt("Name:", "Alf"),
               color: prompt("Color:", "#00ff00"),
               points: 0
            }

            socket.send(JSON.stringify(player)); 
         } else {
            socket.send(JSON.stringify({
               name: "___reconnect___",
               color: document.cookie
            })); 
         }

         socket.onmessage = function (event) {
            if (event.data == 'tick') {
               self.getData();
            }
            else { // stores registered userId
               document.cookie = "#" + event.data;
            }
         };
      };
   }
});

ReactDOM.render(React.createElement(Game, {}), document.getElementById("game"));