const URL = 'https://meet-zode.herokuapp.com'

function FbC() {
    axios.get(`${URL}/getfbc`)
        .then((response) => {
            // Your web app's Firebase configuration
            var firebaseConfig = {
                apiKey: JSON.parse(window.atob(response.data.fbc)).apiKey,
                authDomain: JSON.parse(window.atob(response.data.fbc)).authDomain,
                projectId: JSON.parse(window.atob(response.data.fbc)).projectId,
                storageBucket: JSON.parse(window.atob(response.data.fbc)).storageBucket,
                messagingSenderId: JSON.parse(window.atob(response.data.fbc)).messagingSenderId,
                appId: JSON.parse(window.atob(response.data.fbc)).appId
            };
            // Initialize Firebase
            firebase.initializeApp(firebaseConfig);
        })
}

FbC()

const login = () => {
    const emailDiv = document.getElementById('emailIn')
    const passDiv = document.getElementById('passIn')
    const email = document.getElementById('emailInput').value
    const password = document.getElementById('passwordInput').value

    if (email == "" || password == "") {
        if (email == "") {
            var invalidInput = document.getElementById('invalidEmailInput')
            if (invalidInput == null) {
                var div = document.createElement('div')
                div.className += "invalid-input animate__animated animate__fadeIn animate__faster"
                div.setAttribute('id', 'invalidEmailInput')
                var text = document.createTextNode('Email can\'t be blank')
                div.appendChild(text)
                var element = document.getElementById('emailIn')
                element.appendChild(div)
            }
        }

        if (password == "") {
            var invalidInput = document.getElementById('invalidPasswordInput')
            if (invalidInput == null) {
                var div = document.createElement('div')
                div.className += "invalid-input animate__animated animate__fadeIn animate__faster"
                div.setAttribute('id', 'invalidPasswordInput')
                var text = document.createTextNode('Password can\'t be blank')
                div.appendChild(text)
                var element = document.getElementById('passIn')
                element.appendChild(div)
            }
        }
    } else {
        const loginBtn = document.getElementById('loginBtn')
        loginBtn.innerHTML = `<span id="spinner" class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>&nbsp;&nbsp;Logging in`
        loginBtn.disabled = true
        spinner = document.getElementById('spinner')


        firebase.auth().signInWithEmailAndPassword(email, password)
            .then(() => {
                firebase.auth().currentUser.getIdToken(/* forceRefresh */ true).then((idToken) => {
                    window.location.replace(`${URL}/${MEET_ID}?t=${idToken}`)
                })
            })
            .catch((error) => {
                console.log(error)
                loginBtn.innerHTML = `Login`
                loginBtn.disabled = false
            })
    }
}