// Notification settings
Vue.component('anime', {
    props: ['anime'],
    methods: {
        notify: function() {
            this.anime.subscribed = !this.anime.subscribed;
        }
    },
    template: `
<div class="sub" :class="{ active: anime.subscribed }" @click="notify()">
    <div class="sub-title">{{ anime.title }}</div>
    <div class="sub-status" :class="{ active: anime.subscribed }">{{ anime.subscribed ? "Notify" : "Ignore"}}</div>
</div>`
})

// General settings
Vue.component('set', {
    props: ['option', 'value'],
    methods: {
        set: function(o) {
            this.value.Quality = o;
        },
        sAll: function() {
            this.value.SubscribedAll = !this.value.SubscribedAll
        }
    },
    template: `
<div class="option">
    <div class="option-name">{{ option.name }}</div>
    <div class="option-switch" v-if="option.type == 'choose'">
        <div class="option-button" 
            v-for="(o, i) in option.options" 
            :key="i" 
            :class="{ active: o == value.Quality }"
            @click="set(o)">
            {{ o }}
        </div>
    </div>
    <input class="option-button" 
        v-if="option.type == 'number'"
        v-model="value.Refresh"
        type="number"
        min="1"
        max="3600">
</div>`
})

let app = new Vue({
    el: "#app",
    data: {
        config: {},
        schedule: null,
        illegal: "",

        doneFetching: false
    },
    created: function() {
        fetch("/settings")
        .then(resp => resp.json())
        .then(settings => {
            this.config = settings.config;
            this.schedule = settings.schedule.map(anime => ({
                title: anime.title, 
                time: anime.time, 
                subscribed: !!this.config.Subscriptions.find(x => x == anime.title)}
            )).sort((a, b) => {
                if(a.title < b.title) return -1
                if(a.title > b.title) return 1;
                return 0;
            })
        })
        .then(() => { this.doneFetching = true })
    },
    watch: {
        'config': {
            handler: function(val) {
                if(this.doneFetching) {
                    // Verify refresh
                    let configCopy =  { ...this.config }
                    configCopy.Refresh = parseInt(this.config.Refresh)
                    if(isNaN(configCopy.Refresh) || configCopy.Refresh > 3600 || configCopy.Refresh <= 0 || configCopy.Refresh.length == 0) {
                        configCopy.Refresh = 10
                        this.illegal = "Please enter a number between 1 and 3600";
                    } else {
                        this.illegal = "";
                    }
                    fetch('/settings', { 
                        method: "POST",
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify(configCopy)
                    })
                    .then(resp => {
                        if(!resp.ok) {
                            this.illegal = "Failed to save settings"
                        }
                    })
                }
            },
            deep: true
        },
        'schedule': {
            handler: function(val) {
                if(this.doneFetching) {
                    let notifs = this.schedule.filter(x => x.subscribed);
                    this.config.Subscriptions = notifs.map(anime => anime.title)
                }
            },
            deep: true
        }
    }
})