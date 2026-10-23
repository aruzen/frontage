#include "libfrontage.h"
#include <assert.h>
#include <stdlib.h>
#include <string.h>

static char *copy_string(const char *value)
{
    size_t size = strlen(value) + 1;
    char *copy = malloc(size);
    assert(copy != NULL);
    memcpy(copy, value, size);
    return copy;
}

int main(void)
{
    FreeCardView(NULL);
    FreeBoardView(NULL);

    CardView card = {0};
    FreeCardView(&card);
    card.tag = copy_string("card");
    card.uuid = copy_string("card-id");
    card.cost_count = 2;
    card.costs = calloc(2, sizeof(*card.costs));
    assert(card.costs != NULL);
    FreeCardView(&card);
    assert(card.tag == NULL && card.uuid == NULL && card.costs == NULL);
    assert(card.cost_count == 0);
    FreeCardView(&card);

    BoardView board = {0};
    FreeBoardView(&board);
    board.width = 8;
    board.height = 8;
    board.piece_count = 2;
    board.pieces = calloc(2, sizeof(*board.pieces));
    assert(board.pieces != NULL);
    for (int i = 0; i < board.piece_count; ++i) {
        board.pieces[i].id = copy_string("piece-id");
        board.pieces[i].tag = copy_string("piece");
    }
    board.pieces[1].owner_id = copy_string("owner-id");
    FreeBoardView(&board);
    assert(board.pieces == NULL && board.piece_count == 0);
    assert(board.width == 0 && board.height == 0);
    FreeBoardView(&board);
    return 0;
}
