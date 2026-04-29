document.addEventListener('DOMContentLoaded', () => {
    lucide.createIcons();

    const searchBtn = document.getElementById('search-toggle');
    const searchOverlay = document.getElementById('search-overlay');
    const menuBtn = document.getElementById('menu-toggle');
    const sidebar = document.getElementById('sidebar');
    const sidebarOverlay = document.getElementById('sidebar-overlay');
    const sidebarClose = document.getElementById('sidebar-close');

    if (searchBtn && searchOverlay) {
        searchBtn.addEventListener('click', () => {
            searchOverlay.classList.toggle('active');
            if (searchOverlay.classList.contains('active')) {
                document.querySelector('.search-input')?.focus();
            }
        });
    }

    function toggleSidebar() {
        sidebar?.classList.toggle('active');
        sidebarOverlay?.classList.toggle('active');
    }
    if (menuBtn) menuBtn.addEventListener('click', toggleSidebar);
    if (sidebarClose) sidebarClose.addEventListener('click', toggleSidebar);
    if (sidebarOverlay) sidebarOverlay.addEventListener('click', toggleSidebar);

    const header = document.querySelector('.header');
    window.addEventListener('scroll', () => {
        if (window.scrollY > 10) header?.classList.add('scrolled');
        else header?.classList.remove('scrolled');
    });

    window.Cart = {
        _key: 'clothes_store_cart',

        _read() {
            try { return JSON.parse(localStorage.getItem(this._key)) || []; }
            catch { return []; }
        },
        _write(items) { localStorage.setItem(this._key, JSON.stringify(items)); },

        getAll() { return this._read(); },
        getCount() { return this._read().reduce((s, i) => s + i.qty, 0); },
        getTotal() { return this._read().reduce((s, i) => s + i.qty * i.price, 0); },

        add(product) {
            const items = this._read();
            const key = product.id + '|' + (product.size || '') + '|' + (product.color || '');
            const existing = items.find(i => (i.id + '|' + (i.size || '') + '|' + (i.color || '')) === key);
            if (existing) {
                existing.qty += 1;
            } else {
                items.push({ ...product, qty: 1 });
            }
            this._write(items);
            this.updateBadge();
            Cart.showToast(`${product.name} added to cart`);
        },

        remove(id, size, color) {
            let items = this._read();
            const key = id + '|' + (size || '') + '|' + (color || '');
            items = items.filter(i => (i.id + '|' + (i.size || '') + '|' + (i.color || '')) !== key);
            this._write(items);
            this.updateBadge();
        },

        updateQty(id, size, color, qty) {
            const items = this._read();
            const key = id + '|' + (size || '') + '|' + (color || '');
            const item = items.find(i => (i.id + '|' + (i.size || '') + '|' + (i.color || '')) === key);
            if (item) {
                item.qty = Math.max(1, qty);
                this._write(items);
                this.updateBadge();
            }
        },

        clear() {
            localStorage.removeItem(this._key);
            this.updateBadge();
        },

        updateBadge() {
            const badge = document.getElementById('cart-badge');
            const count = this.getCount();
            if (badge) {
                badge.textContent = count;
                badge.style.display = count > 0 ? 'flex' : 'none';
            }
        },

        showToast(msg) {
            let toast = document.getElementById('store-toast');
            if (!toast) {
                toast = document.createElement('div');
                toast.id = 'store-toast';
                toast.className = 'store-toast';
                document.body.appendChild(toast);
            }
            toast.textContent = msg;
            toast.classList.add('show');
            setTimeout(() => toast.classList.remove('show'), 2200);
        }
    };

    window.Wishlist = {
        _key: 'clothes_store_wishlist',

        _read() {
            try { return JSON.parse(localStorage.getItem(this._key)) || []; }
            catch { return []; }
        },
        _write(items) { localStorage.setItem(this._key, JSON.stringify(items)); },

        getAll() { return this._read(); },
        has(id) { return this._read().some(i => i.id === id); },

        toggle(product) {
            let items = this._read();
            const idx = items.findIndex(i => i.id === product.id);
            if (idx >= 0) {
                items.splice(idx, 1);
                Cart.showToast(`${product.name} removed from wishlist`);
            } else {
                items.push(product);
                Cart.showToast(`${product.name} added to wishlist`);
            }
            this._write(items);
        },

        remove(id) {
            let items = this._read();
            items = items.filter(i => i.id !== id);
            this._write(items);
        }
    };

    Cart.updateBadge();

    function requireSizeAndColor(size, color, message) {
        if (!size || !color) {
            Cart.showToast(message || 'Select size and color');
            return false;
        }
        return true;
    }

    document.querySelectorAll('.product-card[data-product-id]').forEach(card => {
        const id = card.dataset.productId;
        const name = card.dataset.productName;
        const price = parseFloat(card.dataset.productPrice);
        const image = card.dataset.productImage || '';

        const cartBtn = card.querySelector('.btn-add-cart');
        const wishBtn = card.querySelector('.btn-add-wish');

        if (cartBtn) {
            cartBtn.addEventListener('click', e => {
                e.preventDefault();
                e.stopPropagation();
                if (!requireSizeAndColor('', '', 'Select size and color on the product page')) {
                    return;
                }
                Cart.add({ id, name, price, image });
            });
        }
        if (wishBtn) {
            wishBtn.addEventListener('click', e => {
                e.preventDefault();
                e.stopPropagation();
                Wishlist.toggle({ id, name, price, image });
            });
        }
    });

    const productAddCartBtn = document.getElementById('product-add-cart');
    if (productAddCartBtn) {
        productAddCartBtn.addEventListener('click', e => {
            e.preventDefault();
            const id = productAddCartBtn.dataset.productId;
            const name = productAddCartBtn.dataset.productName;
            const price = parseFloat(productAddCartBtn.dataset.productPrice);
            const image = productAddCartBtn.dataset.productImage || '';
            const size = document.querySelector('input[name="size"]:checked')?.value || '';
            const color = document.querySelector('input[name="color"]:checked')?.value || '';
            if (!requireSizeAndColor(size, color)) {
                return;
            }
            Cart.add({ id, name, price, image, size, color });
        });
    }

    const productWishBtn = document.getElementById('product-add-wish');
    if (productWishBtn) {
        productWishBtn.addEventListener('click', () => {
            const id = productWishBtn.dataset.productId;
            const name = productWishBtn.dataset.productName;
            const price = parseFloat(productWishBtn.dataset.productPrice);
            const image = productWishBtn.dataset.productImage || '';
            Wishlist.toggle({ id, name, price, image });
        });
    }

    const wishlistContainer = document.getElementById('wishlist-items');
    const wishlistEmpty = document.getElementById('wishlist-empty');

    if (wishlistContainer) {
        renderWishlist();
    }

    function renderWishlist() {
        const items = Wishlist.getAll();
        wishlistContainer.innerHTML = '';

        if (items.length === 0) {
            wishlistContainer.style.display = 'none';
            if (wishlistEmpty) wishlistEmpty.style.display = 'block';
            return;
        }

        wishlistContainer.style.display = 'grid';
        if (wishlistEmpty) wishlistEmpty.style.display = 'none';

        items.forEach(item => {
            const card = document.createElement('div');
            card.className = 'product-card';
            card.innerHTML = `
                <a href="/product/${item.id}">
                    <div class="product-image-container">
                        <img src="${item.image}" alt="${item.name}" class="product-image"
                             onerror="this.parentElement.style.background='#f5f5f5'">
                        <div class="product-actions" style="opacity:1;transform:none;">
                            <button class="action-btn btn-wish-cart" title="Add to Cart">
                                <i data-lucide="shopping-bag" size="18"></i>
                            </button>
                            <button class="action-btn btn-wish-remove" title="Remove"
                                    style="background:var(--color-danger);color:white;">
                                <i data-lucide="trash-2" size="18"></i>
                            </button>
                        </div>
                    </div>
                    <div class="product-info">
                        <div>
                            <div class="product-title">${item.name}</div>
                            <div class="product-price">$${parseFloat(item.price).toFixed(2)}</div>
                        </div>
                    </div>
                </a>
            `;

            card.querySelector('.btn-wish-cart').addEventListener('click', e => {
                e.preventDefault();
                e.stopPropagation();
                if (!requireSizeAndColor('', '', 'Select size and color on the product page')) {
                    return;
                }
                Cart.add({ id: item.id, name: item.name, price: item.price, image: item.image });
            });

            card.querySelector('.btn-wish-remove').addEventListener('click', e => {
                e.preventDefault();
                e.stopPropagation();
                Wishlist.remove(item.id);
                renderWishlist();
                lucide.createIcons();
            });

            wishlistContainer.appendChild(card);
        });

        lucide.createIcons();
    }

    const cartItemsEl = document.getElementById('cart-items');
    const cartEmptyEl = document.getElementById('cart-empty');
    const cartSummaryEl = document.getElementById('cart-summary');

    if (cartItemsEl) {
        renderCart();
    }

    function renderCart() {
        const items = Cart.getAll();
        cartItemsEl.innerHTML = '';

        if (items.length === 0) {
            cartItemsEl.style.display = 'none';
            if (cartSummaryEl) cartSummaryEl.style.display = 'none';
            if (cartEmptyEl) cartEmptyEl.style.display = 'flex';
            return;
        }

        cartItemsEl.style.display = 'block';
        if (cartSummaryEl) cartSummaryEl.style.display = 'block';
        if (cartEmptyEl) cartEmptyEl.style.display = 'none';

        items.forEach(item => {
            const row = document.createElement('div');
            row.className = 'cart-item';
            const variant = [item.size, item.color].filter(Boolean).join(' / ') || 'Standard';
            row.innerHTML = `
                <div class="cart-item-image">
                    <img src="${item.image}" alt="${item.name}"
                         onerror="this.parentElement.style.background='#f5f5f5'">
                </div>
                <div class="cart-item-details">
                    <a href="/product/${item.id}" class="cart-item-name">${item.name}</a>
                    <div class="cart-item-variant">${variant}</div>
                    <div class="cart-item-price">$${parseFloat(item.price).toFixed(2)}</div>
                </div>
                <div class="cart-item-qty">
                    <button class="qty-btn qty-minus">−</button>
                    <span class="qty-value">${item.qty}</span>
                    <button class="qty-btn qty-plus">+</button>
                </div>
                <div class="cart-item-total">$${(item.price * item.qty).toFixed(2)}</div>
                <button class="cart-item-remove" title="Remove">
                    <i data-lucide="x" size="18"></i>
                </button>
            `;

            row.querySelector('.qty-minus').addEventListener('click', () => {
                if (item.qty > 1) {
                    Cart.updateQty(item.id, item.size, item.color, item.qty - 1);
                    renderCart();
                    lucide.createIcons();
                }
            });
            row.querySelector('.qty-plus').addEventListener('click', () => {
                Cart.updateQty(item.id, item.size, item.color, item.qty + 1);
                renderCart();
                lucide.createIcons();
            });
            row.querySelector('.cart-item-remove').addEventListener('click', () => {
                Cart.remove(item.id, item.size, item.color);
                renderCart();
                lucide.createIcons();
            });

            cartItemsEl.appendChild(row);
        });

        const subtotal = Cart.getTotal();
        const delivery = subtotal > 0 ? 20 : 0;
        const total = subtotal + delivery;

        document.getElementById('cart-subtotal').textContent = '$' + subtotal.toFixed(2);
        document.getElementById('cart-delivery').textContent = subtotal > 200 ? 'Free' : '$' + delivery.toFixed(2);
        document.getElementById('cart-total').textContent = '$' + (subtotal > 200 ? subtotal : total).toFixed(2);
        document.getElementById('cart-count-label').textContent = Cart.getCount() + ' item(s)';

        lucide.createIcons();
    }
});
