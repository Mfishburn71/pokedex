// Boot the WASM runtime
const go = new Go();
WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject)
    .then(result => {
        go.run(result.instance);
    });

// Called by Go when no save exists
function showNameModal() {
    document.getElementById("name-modal").classList.add("active");
}


function prettify(name) {
    return name.split("-").map(w => w.charAt(0).toUpperCase() + w.slice(1)).join(" ");
}

// Called by Go after loading a save
function renderTrainerHeader(name) {
    document.getElementById("right-content").textContent = "Trainer: " + name;
}

document.getElementById("trainer-name-submit").addEventListener("click", () => {
    const name = document.getElementById("trainer-name-input").value.trim();
    if (!name) return;
    cfg_setName(name);  // we'll add this Go wrapper later
    document.getElementById("name-modal").classList.remove("active");
    renderTrainerHeader(name);
});

// Utility: render text to the left screen
function renderLeft(html, clearRight = true) {
    const content = document.getElementById("left-content");
    const indicator = document.getElementById("scroll-indicator");
       if (clearRight) {
        document.getElementById("right-content").innerHTML = "";
    }
    content.innerHTML = html;

    setTimeout(() => {
        const hasOverflow = content.scrollHeight > content.clientHeight;
        console.log("scrollHeight:", content.scrollHeight, "clientHeight:", content.clientHeight, "hasOverflow:", hasOverflow);
        console.log("Setting indicator to:", hasOverflow ? "block" : "none");
        indicator.style.display = hasOverflow ? "block" : "none";

        content.onscroll = () => {
            const atBottom = content.scrollTop + content.clientHeight >= content.scrollHeight - 4;
            indicator.style.display = atBottom ? "none" : "block";
        };
    }, 0);
}

// Utility: render text to the right screen
function renderRight(text) {
    document.getElementById("right-content").textContent = text;
}

// Button wiring
document.getElementById("btn-help").addEventListener("click", () => {
    const result = JSON.parse(showHelp());
    const html = result.map(e => `<p><strong>${e.name}</strong>: ${e.description}</p>`).join("");
    renderLeft(html);
});

document.getElementById("btn-map").addEventListener("click", () => {
    const result = JSON.parse(listMap());
    const html = result.map(loc => `<p>${prettify(loc.name)}</p>`).join("");
    renderLeft(html);
});

document.getElementById("btn-mapb").addEventListener("click", () => {
    const raw = listMapB();
    const result = JSON.parse(raw);
    if (result.error) {
        renderRight(prettify(result.error));
        return;
    }
    const html = result.map(loc => `<p>${prettify(loc.name)}</p>`).join("");
    renderLeft(html);
});


document.getElementById("btn-catch").addEventListener("click", () => {
    const name = prompt("Which Pokemon?");
    if (!name) return;

    const raw = catchPokemon(name.toLowerCase().trim());
    const result = JSON.parse(raw);

    if (result.error || result.message.includes("already in your Pokedex")) {
        renderRight(result.error || result.message);
        return;
    }

    // Show the target sprite before throwing
    if (result.sprite) {
        renderLeft(`<img src="${result.sprite}" alt="${result.name}" style="width: 100%; height: 100%; object-fit: contain;">`, false);
    }

    const ball = document.createElement("img");
    ball.src = "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/items/poke-ball.png";
    ball.classList.add("pokeball-throw");
    document.querySelector(".screen-right").appendChild(ball);

    ball.addEventListener("animationend", () => {
    ball.remove();
    renderRight(result.message);
    if (!result.caught) {
        renderLeft("", false);  // don't clear right screen
    }
});
});

document.getElementById("btn-inspect").addEventListener("click", () => {
    const name = prompt("Inspect which Pokemon?");
    if (!name) return;
    const result = JSON.parse(inspectPokemon(name.toLowerCase().trim()));
    if (result.error) {
        renderRight(prettify(result.error));
        return;
    }
    const statsHtml = result.stats.map(s => `<p>${s.name}: ${s.value}</p>`).join("");
    renderLeft(`
        <img src="${result.sprite}" alt="${result.display}">
        <p><strong>${result.display}</strong></p>
        <p>Types: ${result.types.join(", ")}</p>
        ${statsHtml}
    `);
    document.getElementById("right-content").innerHTML = 
    `${result.display} | Caught in: ${prettify(result.ball_type)}<br>
    Height: ${result.height} | Weight: ${result.weight}<br>
    Current Level: 1 | Affection: Normal
    `;
    
});

document.getElementById("btn-ball").addEventListener("click", () => {
    const ball = prompt("The following balls are in your bag right now: Pokeball, Great Ball, Ultra Ball, Master Ball");
    if (!ball) return;
    const result = JSON.parse(equipokeball(ball.toLowerCase().trim()));
    if (result.error) {
        renderRight(result.error);
        return;
    }
    renderRight(result.message);
});

document.getElementById("btn-battle").addEventListener("click", () => {
    const name1 = prompt("First Pokemon?");
    if (!name1) return;
    const name2 = prompt("Second Pokemon?");
    if (!name2) return;
    const result = JSON.parse(createBattle(name1.toLowerCase().trim(), name2.toLowerCase().trim()));
    if (result.error) {
        renderRight(result.error);
        return;
    }
    const html = result.logs.map(log => `<p>${log}</p>`).join("");
    renderLeft(html);

    // Get sprites from pokedex
    const pokedex = JSON.parse(openPokedex());
    const mon1 = pokedex[name1.toLowerCase().trim()];
    const mon2 = pokedex[name2.toLowerCase().trim()];
    const sprite1 = mon1?.pokemon?.sprites?.front_default || "";
    const sprite2 = mon2?.pokemon?.sprites?.front_default || "";

    document.getElementById("right-content").innerHTML = `
        <div style="display:flex; justify-content:space-around; align-items:center; height:100%">
            <img src="${sprite1}" alt="${name1}" id="battle-sprite-1">
            <img src="${sprite2}" alt="${name2}" id="battle-sprite-2">
        </div>
    `;

    // Highlight the winner
    const lastLog = result.logs[result.logs.length - 1];
  if (lastLog.includes(result.fighter1)) {
    document.getElementById("battle-sprite-1").style.filter = "drop-shadow(0 0 8px gold)";
    document.getElementById("battle-sprite-2").style.filter = "grayscale(100%)";
} else if (lastLog.includes(result.fighter2)) {
    document.getElementById("battle-sprite-2").style.filter = "drop-shadow(0 0 8px gold)";
    document.getElementById("battle-sprite-1").style.filter = "grayscale(100%)";
}
});

document.getElementById("btn-party").addEventListener("click", () => {
    const result = JSON.parse(listParty());
    if (result.error) {
        renderRight(result.error);
        return;
    }
    
    renderParty(result.party);
});

document.getElementById("btn-explore").addEventListener("click", () => {
    const area = prompt("Explore which area?");
    if (!area) return;
    const normalized = area.toLowerCase().trim().replace(/\s+/g, "-");
    const result = JSON.parse(exploreArea(normalized));
    if (result.error) {
        renderRight(result.error);
        return;
    }
    const pokemonHtml = result.pokemon.map(p => `<p>${p.display} (${p.methods.join(", ")})</p>`).join("");
    renderLeft(`
        <p><strong>${result.display}</strong></p>
        <p>Tier: ${result.tier} | Avg Level: ${result.level}</p>
        <p>Regions: ${result.regions.join(", ")}</p>
        <p>Habitats: ${result.habitats.join(", ")}</p>
        ${pokemonHtml}
    `);
});

document.getElementById("btn-search").addEventListener("click", () => {
    const query = prompt("Search for area:");
    if (!query) return;
    renderRight(`Waking up the professor. Please be patient!`);
    setTimeout(() => { //Pauses to show the wait message
        const result = JSON.parse(searchArea(query.toLowerCase().trim()));
        if (result.error) {
            renderRight(result.error);
            return;
        }
        const html = result.results.map(r => `<p>${r}</p>`).join("");
        renderLeft(`<p>Found ${result.count} results:</p>${html}`);
    }, 0);
});


document.getElementById("btn-region").addEventListener("click", () => {
    
    const region = prompt("Which region?");
    if (!region) return;
    const query = prompt("Filter by name (leave blank for all):");
    renderRight(`Waking up the professor. Please be patient!`);
    setTimeout(() => { //Pauses to show the wait message
    const result = JSON.parse(searchRegion(region.toLowerCase().trim(), (query || "").toLowerCase().trim()));
    if (result.error) {
        renderRight(prettify(result.error));
        return;
    }
    const html = result.results.map(r => `<p>${prettify(r)}</p>`).join("");
    renderLeft(`<p>${prettify(result.region)} - ${result.count} Areas:</p>${html}`);
}, 2);
});

document.getElementById("btn-pokedex").addEventListener("click", () => {
    const result = JSON.parse(openPokedex());
    if (result.error) {
        renderRight(result.error);
        return;
    }
    const entries = Object.entries(result);
    if (entries.length === 0) {
        renderLeft("<p>Your Pokedex is empty!</p>");
        renderRight("Total Caught: 0");
        return;
    }
    const html = entries.map(([key, val]) => {
        const displayName = val.nickname 
            ? `${prettify(val.nickname)} <em>(${prettify(val.pokemon.name)})</em>` 
            : prettify(val.pokemon.name);
        return `
            <p>
                <strong>${displayName}</strong>
                <button onclick="handleNickname('${key}')">Nickname</button>
                <button onclick="handlePartyAdd('${key}')">+ Party</button>
                <button onclick="handlePartyRemove('${key}')">- Party</button>
            </p>
        `;
    }).join("");
    renderLeft(html);

const info = JSON.parse(getTrainerInfo());
console.log("trainer info:", info);
document.getElementById("right-content").innerHTML = 
    `Pokedex registered to: Trainer ${info.name}<br>In your hand: ${prettify(info.ball)}<br>Total Pokemon Caught: ${entries.length}<br>Emblems (Badges): 0`;
});

document.getElementById("btn-data").addEventListener("click", () => {
    const action = prompt("Type: save, load, or delete");
    if (!action) return;
    switch (action.toLowerCase().trim()) {
        case "save": {
            const data = saveData();
            localStorage.setItem("pokedex_save", data);
            renderRight("Game saved!");
            break;
        }
        case "load": {
            const saved = localStorage.getItem("pokedex_save");
            if (!saved) {
                renderRight("No save data found.");
                return;
            }
            const result = JSON.parse(loadData(saved));
            if (result.error) {
                renderRight(result.error);
                return;
            }
            renderRight(result.message);
            break;
        }
        case "delete": {
            const confirm = prompt("Are you sure? Type YES to confirm.");
            if (confirm !== "YES") return;
            const result = JSON.parse(deleteData());
            localStorage.removeItem("pokedex_save");
            renderRight(result.message);
            break;
        }
        default:
            renderRight("Unknown action.");
    }
});

// Contextual handlers called from inline buttons in the left screen

function handleNickname(pokemonName) {
    const nickname = prompt(`Nickname for ${pokemonName}?`);
    if (!nickname) return;
    const result = JSON.parse(setNickname(pokemonName, nickname.trim()));
    if (result.error) {
        renderRight(result.error);
        return;
    }
    renderRight(result.message);
    // Refresh the pokedex view
    document.getElementById("btn-pokedex").click();
}

function handlePartyAdd(pokemonName) {
    const result = JSON.parse(addParty(pokemonName));
    if (result.error) {
        renderRight(result.error);
        return;
    }
    renderRight(result.message);
    renderParty(result.party);
}

function handlePartyRemove(pokemonName) {
    const result = JSON.parse(removeParty(pokemonName));
    if (result.error) {
        renderRight(result.error);
        return;
    }
    renderRight(result.message);
    renderParty(result.party);
}
function renderParty(party) {
    if (!party || party.length === 0) {
        renderLeft("<p>Your party is empty.</p>");
        document.getElementById("right-content").innerHTML = "";
        return;
    }
    const html = party.map(m => `
        <p>Slot ${m.slot}: <strong>${m.display}</strong>
            <button onclick="handlePartyRemove('${m.name}')">Remove</button>
        </p>
    `).join("");
    renderLeft(html);

    const pokedex = JSON.parse(openPokedex());

    let maxH = party.length > 3 ? "70%" : "100%";
    let maxW = party.length > 3 ? "65%" : "65%";
    maxH = party.length > 4 ? "55%" : maxH;
    maxW = party.length > 4 ? "55%" : maxW;
    maxH = party.length > 5 ? "47%" : maxH;
    maxW = party.length > 5 ? "50%" : maxW;
     
    const sprites = party
    .map(m => pokedex[m.name])
    .filter(p => p && p.pokemon?.sprites?.front_default)

    .map(p => `<img src="${p.pokemon.sprites.front_default}" alt="${p.pokemon.name}" style="max-height:${maxH}; max-width:${maxW}; object-fit:contain;">`)

    //.map(p => `<img src="${p.pokemon.sprites.front_default}" alt="${p.pokemon.name}" style="max-height:${maxH}; width:${maxW}; object-fit:contain;">`)
    .join("");

    document.getElementById("right-content").innerHTML = `
        <div style="display:flex; flex-wrap:wrap; justify-content:center; align-items:center; height:100%">
            ${sprites}
        </div>
    `;
}

// cfg_setName - we need a Go wrapper for this
// For now, store the name in the save data on first save
function cfg_setName(name) {
    setTrainerName(name);
}

function syncFetch(url) {
    try {
        const xhr = new XMLHttpRequest();
        xhr.open("GET", url, false);
        xhr.send();
        if (xhr.status > 299) {
            return { error: "bad status code: " + xhr.status };
        }
        const encoder = new TextEncoder();
        const encoded = encoder.encode(xhr.responseText);
        console.log("syncFetch response length:", xhr.responseText.length);
        console.log("syncFetch first 100 chars:", xhr.responseText.substring(0, 100));
        return { data: encoded };
    } catch (e) {
        return { error: e.message };
    }
}
