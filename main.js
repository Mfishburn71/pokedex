const help = JSON.parse(showHelp());
help.forEach(cmd => {
    // render each entry however you like
});

const result = JSON.parse(catchPokemon("pikachu"));
if (result.error) {
    showMessage(result.error); // your existing error text appears here
} else {
    renderCatch(result);
}

const save = localStorage.getItem("pokedex_save");
if (save) {
    const result = loadSave(save); // Go function
    renderTrainerHeader(JSON.parse(result));
} else {
    showNameModal();
}

const result = JSON.parse(createBattle("pikachu", "squirtle"));
if (result.error) {
    showMessage(result.error);
    return;
}

let i = 0;
const interval = setInterval(() => {
    if (i >= result.logs.length) {
        clearInterval(interval);
        return;
    }
    appendBattleLog(result.logs[i]);
    i++;
}, 800); // one line every 800ms


function exitPokedex() {
    // save first
    const saveJSON = saveData();
    localStorage.setItem("pokedex_save", saveJSON);
    // then do whatever makes sense in your UI
    // show start screen, display goodbye message, etc.
}

// saving
const saveJSON = saveData();
localStorage.setItem("pokedex_save", saveJSON);

// loading
const saved = localStorage.getItem("pokedex_save");
if (saved) loadData(saved);

// deleting
deleteData();
localStorage.removeItem("pokedex_save");